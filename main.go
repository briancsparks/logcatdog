package main

import (
  "bcs/logcatdog/kafka"
  "bcs/logcatdog/proc"
  "bcs/logcatdog/tries"
  "fmt"
  "log"
  "sync"
)

func main() {
  tries.M()
}

func mainX() {
	fmt.Printf("%s\n", "booya")

	var wg sync.WaitGroup

  //// Trap SIGINT to trigger a graceful shutdown.
  //signals := make(chan os.Signal, 1)
  //signal.Notify(signals, os.Interrupt)

  adb := proc.NewAdb()

  adbDone, err := adb.Initit()
  if err != nil {
    log.Fatal(err)
  }

  producer := kafka.NewProducer(adb.DataStream)

  wg.Add(1)
  go func() {
    defer wg.Done()

    pDone, err := producer.Send()
    if err != nil {
      log.Fatal(err)
    }

    <- pDone
  }()

  wg.Add(1)
  go func() {
    defer wg.Done()

    for d := range adb.DataStream2 {
      if d.LineNum >= 100 {
        adb.Stop()
      }

      if d.LineNum <= 100 {
        fmt.Printf("linebyline %d|%s|\n", d.LineNum, d.Line)
      }
    }

    adb.Term(0)

    <-adbDone
  }()

  adb.Start <- struct{}{}

  wg.Wait()

  fmt.Printf("%s\n", "booya too")
}



