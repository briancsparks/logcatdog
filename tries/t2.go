package tries

import (
  "bcs/logcatdog/kafka"
  lbl "bcs/logcatdog/linebyline"
  "bcs/logcatdog/proc"
  "fmt"
  "log"
  "sync"
)

func M() {

  // Get list of devices
  devices, err := proc.AdbDevices()
  if err != nil {
    log.Fatal(err)
  }

  for n, device := range devices {
    fmt.Printf("device %d: %v\n", n, device)
  }

  // `done` for pipeline
  done := make(chan struct{})
  //defer close(done)
  closed := false
  closeIt := func() {
    if !closed {
      closed = true
      close(done)
    }
  }
  defer closeIt()

  // adb logcat ...
  adbLines, err := proc.Launch(done, "adb", "-d", "logcat", "-d", "-b", "all", "-v", "threadtime", "-v", "usec", "-v", "year", "-v", "UTC", "-v", "epoch")
  if err != nil {
    log.Fatal(err)
  }

  adbLines1, adbLines2 := lbl.Dup(adbLines)

  wg := sync.WaitGroup{}

  wg.Add(1)
  go func() {
    defer wg.Done()

    // Read the stream of lines
    for d := range adbLines1 {
      if d.LineNum >= 100 {
        closeIt()
      }

      fmt.Printf("linebyline %d|%s|\n", d.LineNum, d.Line)
    }
  }()

  kafka.Send(done, adbLines2, "logcatter")

  wg.Wait()

}
