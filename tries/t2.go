package tries

import (
  "bcs/logcatdog/kafka"
  lbl "bcs/logcatdog/linebyline"
  "bcs/logcatdog/proc"
  "fmt"
  "log"
  "regexp"
  "sync"
  "time"
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

  m1 := regexp.MustCompile(`[^a-zA-Z0-9]`)
  date := m1.ReplaceAllString(time.Now().UTC().Format("2006-01-02T15:04:05.000000000Z"), "_")
  kafka.Send(done, adbLines2, "logcatter_" + date)

  wg.Wait()

}
