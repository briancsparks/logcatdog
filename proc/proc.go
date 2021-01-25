package proc

import (
  "log"
  "os/exec"
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

