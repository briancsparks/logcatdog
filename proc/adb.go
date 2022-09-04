package proc

import (
  "bufio"
  "fmt"
  "log"
  "os/exec"
  "regexp"
  "strings"

  "bcs/logcatdog/linebyline"
)


type Adb struct {

  DataStream chan linebyline.LineData
  DataStream2 chan linebyline.LineData

  Start chan struct{}
  stop  chan struct{}
  done  chan struct{}

  // Windows
  peg ProcessExitGroup
}

func NewAdb() Adb {
  adb := new(Adb)

  adb.DataStream = make(chan linebyline.LineData)
  adb.DataStream2 = make(chan linebyline.LineData)
  adb.Start = make(chan struct{})
  adb.stop = make(chan struct{})
  adb.done = make(chan struct{})

  return *adb
}

func (adb *Adb) Stop() {
  go func() {
    adb.stop <- struct{}{}
  }()
}

func (adb *Adb) Term(exitcode uint32) {
  adb.peg.Terminate(exitcode)
}

func (adb *Adb) Initit() (chan struct{}, error) {

  go func() {
    <- adb.Start

    devices, err := AdbDevices()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("devices: %v\n", devices)

    cmd := exec.Command("adb", "-d", "logcat", "-d", "-b", "all", "-v", "threadtime", "-v", "usec", "-v", "year", "-v", "UTC", "-v", "epoch")

    // TODO: Kill on non-Windows
    // Killing a sub-process: https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
    //cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    //syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

    adbStdout, err := cmd.StdoutPipe()
    if err != nil {
      log.Fatal(err)
    }

    scanner := bufio.NewScanner(adbStdout)
    go func() {
      //defer close(adb.DataStream)
      //defer close(adb.DataStream2)

      count := 0
      //fmt.Printf("entering scan0()\n")
scanLoop:
      for scanner.Scan() {
        //fmt.Printf("got item from scan()\n")
        //fmt.Printf("ADBBBBB |%s\n", scanner.Text())
        l := linebyline.NewLineData(count, scanner.Text(), "-")
        adb.DataStream <- l
        adb.DataStream2 <- l
        count += 1

        // stop?
        //fmt.Printf("selecting\n")
        select {
        case <-adb.stop:
          //fmt.Printf("select adb.stopping\n")
          break scanLoop
        default:
          // nothing, not stopping
          //fmt.Printf("select defaulting\n")
        }
        //fmt.Printf("entering scan()\n")
      }
      //fmt.Printf("exiting scan()\n")

      close(adb.DataStream)
      close(adb.DataStream2)
    }()

    if err = cmd.Start(); err != nil {
      log.Fatal(err)
    }

    peg, err := NewProcessExitGroup()
    if err != nil {
      log.Fatal(err)
    }

    adb.peg = peg
    if err = adb.peg.AddProcess(cmd.Process); err != nil {
      log.Fatal(err)
    }

    if err = cmd.Wait(); err != nil {
      log.Fatal(err)
    }

    adb.peg.Dispose()
    //if err = adb.peg.Dispose(); err != nil {
    //  log.Fatal(err)
    //}

    adb.done <- struct{}{}
  }()

  return adb.done, nil
}

// AdbDevices --------------------------------------------------------------------------------------------------------------------
func AdbDevices() ([]string, error) {
  output, err := QuickLaunch("adb", "-d", "devices")
  if err != nil {
    log.Fatal(err)
  }

  re := regexp.MustCompile(`^(\S+)\s+device`)

  devices := make([]string, 0)
  for _, line := range strings.Split(output, "\n") {
    if parts := re.FindAllStringSubmatch(line, 1); parts != nil {
      devices = append(devices, parts[0][1])
    }
  }

  return devices, nil
}

