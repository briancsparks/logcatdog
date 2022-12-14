package proc

import (
  "os"
  "unsafe"

  "golang.org/x/sys/windows"
)

// Essentially from: https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209
// Also good: https://github.com/alexbrainman/ps

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type process struct {
  Pid    int
  Handle uintptr
}

type ProcessExitGroup windows.Handle

func NewProcessExitGroup() (ProcessExitGroup, error) {
  handle, err := windows.CreateJobObject(nil, nil)
  if err != nil {
    return 0, err
  }

  info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
    BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
      LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
    },
  }
  if _, err := windows.SetInformationJobObject(
    handle,
    windows.JobObjectExtendedLimitInformation,
    uintptr(unsafe.Pointer(&info)),
    uint32(unsafe.Sizeof(info))); err != nil {
    return 0, err
  }

  return ProcessExitGroup(handle), nil
}

func (g ProcessExitGroup) AddProcess(p *os.Process) error {
  return windows.AssignProcessToJobObject(
    windows.Handle(g),
    windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}

func (g ProcessExitGroup) Terminate(exitcode uint32)  {
  _ = windows.TerminateJobObject(windows.Handle(g), exitcode)
}

func (g ProcessExitGroup) Terminate2(exitcode uint32) error {
  return windows.TerminateJobObject(windows.Handle(g), exitcode)
}

func (g ProcessExitGroup) Dispose() {
  _ = windows.CloseHandle(windows.Handle(g))
}

func (g ProcessExitGroup) Dispose2() error {
  return windows.CloseHandle(windows.Handle(g))
}

