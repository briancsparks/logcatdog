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
      fmt.Printf("ADBBBBB |%d|%s\n", d.LineNum, d.Line)
      if d.LineNum >= 100 {
        adb.Stop()
      }
    }

    adb.Term(0)
  }()

  adb.Start <- struct{}{}

  <-done

  fmt.Printf("%s\n", "booya too")
}


//func goadb() {
//
//  devices, err := adbDevices()
//  if err != nil {
//    log.Fatal(err)
//  }
//  fmt.Printf("devices: %v\n", devices)
//
//
//  // Killing a sub-process: https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
//
//  cmd := exec.Command("adb", "-d", "logcat" ,"-d", "-b", "all", "-v", "threadtime", "-v", "usec", "-v", "year", "-v", "UTC", "-v", "epoch")
//  //cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
//  //syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
//  adbStdout, err := cmd.StdoutPipe()
//  if err != nil {
//    log.Fatal(err)
//  }
//
//  scanner := bufio.NewScanner(adbStdout)
//  go func() {
//    for scanner.Scan() {
//      fmt.Printf("ADBBBBB |%s\n", scanner.Text())
//    }
//  }()
//
//  if err = cmd.Start(); err != nil {
//    log.Fatal(err)
//  }
//
//  if err = cmd.Wait(); err != nil {
//    log.Fatal(err)
//  }
//}


