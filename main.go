package main

import (
  "bcs/logcatdog/proc"
  "bufio"
  "fmt"
  "log"
  "os/exec"
  "regexp"
  "strings"
)

func main() {
	fmt.Printf("%s\n", "booya")

	goadb()

  fmt.Printf("%s\n", "booya too")
}

func goadb() {

  devices, err := adbDevices()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("devices: %v\n", devices)

  cmd := exec.Command("adb", "-d", "logcat" ,"-d", "-b", "all", "-v", "threadtime", "-v", "usec", "-v", "year", "-v", "UTC", "-v", "epoch")
  adbStdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }

  scanner := bufio.NewScanner(adbStdout)
  go func() {
    for scanner.Scan() {
      fmt.Printf("ADBBBBB |%s\n", scanner.Text())
    }
  }()

  if err = cmd.Start(); err != nil {
    log.Fatal(err)
  }

  if err = cmd.Wait(); err != nil {
    log.Fatal(err)
  }
}

func adbDevices() ([]string, error) {
  output, err := proc.QuickLaunch("adb", "-d", "devices")
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


