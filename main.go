package main

import (
  "bcs/logcatdog/proc"
  "fmt"
  "log"
)

func main() {
	fmt.Printf("%s\n", "booya")

  adb := proc.NewAdb()

  done, err := adb.Initit()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    for d := range adb.DataStream {
      if d.LineNum >= 100 {
        adb.Stop()
      }

      if d.LineNum <= 100 {
        fmt.Printf("line %d|%s|\n", d.LineNum, d.Line)
      }
    }

    adb.Term(0)
  }()

  adb.Start <- struct{}{}

  <-done

  fmt.Printf("%s\n", "booya too")
}



