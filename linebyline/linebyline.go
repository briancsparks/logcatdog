package linebyline

//import "fmt"

type LineData struct {
  Filename string `json:"filename"`
  LineNum int     `json:"num"`
  Line string     `json:"line"`
}

func NewLineData(lineNum int, line, filename string) LineData {
  l := new(LineData)

  l.Filename = filename
  l.Line = line
  l.LineNum = lineNum

  return *l
}

func Dup(l chan LineData) (chan LineData, chan LineData) {
  a := make(chan LineData)
  b := make(chan LineData)

  go func() {
    defer close(a)
    defer close(b)

    for line := range l {

      func (line LineData) {
        a <- line
        b <- line
      }(line)
    }
  }()

  return a, b
}

func DupAbortable(done <-chan struct{}, l chan LineData) (chan LineData, chan LineData) {
  a := make(chan LineData)
  b := make(chan LineData)

  go func() {
    defer close(a)
    defer close(b)

    for line := range l {

      func (line LineData) {
        select {
        case a <- line:
          b <- line

        case <- done:
          return
        }
      }(line)
    }
  }()

  return a, b
}



