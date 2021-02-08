package proc

import (
  "bufio"
  "log"
  "os/exec"

  // Ours
  lbl "bcs/logcatdog/linebyline"
)

func QuickLaunch(name string, args ...string) (string, error) {
  var (
    output []byte
    err    error
  )

  if output, err = exec.Command(name, args...).Output(); err != nil {
    log.Fatal(err)
  }

  result := string(output)
  return result, nil
}

func Launch(done <-chan struct{}, name string, args ...string) (chan lbl.LineData, error) {
  out := make(chan lbl.LineData)

  go func() {
    defer close(out)

    //cmd := exec.Command("adb", "-d", "logcat", "-d", "-b", "all", "-v", "threadtime", "-v", "usec", "-v", "year", "-v", "UTC", "-v", "epoch")
    cmd := exec.Command(name, args...)

    outPipe, err := cmd.StdoutPipe()
    if err != nil {
      log.Fatal(err)
    }

    peg, err := NewProcessExitGroup()
    if err != nil {
      log.Fatal(err)
    }
    defer peg.Dispose()

    scanner := bufio.NewScanner(outPipe)

    go func() {
      defer close(out)

      count := 0
      for scanner.Scan() {

        select {
        case out <- lbl.NewLineData(count, scanner.Text(), "-"):
          count += 1

        case <- done:
          peg.Terminate(8)
          return
        }
      }

    }()

    if err = cmd.Start(); err != nil {
      log.Fatal(err)
    }

    if err = peg.AddProcess(cmd.Process); err != nil {
      log.Fatal(err)
    }

    if err = cmd.Wait(); err != nil {
      log.Fatal(err)
    }

    //if err = peg.Dispose(); err != nil {
    //  log.Fatal(err)
    //}

  }()

  return out, nil
}

