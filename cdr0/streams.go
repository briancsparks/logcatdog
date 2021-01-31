package cdr0

type Streams struct {
  items map[string]chan string
}

func NewStreams() Streams {
  s := new(Streams)
  s.items = make(map[string]chan string)
  return *s
}

func NewStreams1(name string, stream chan string) Streams {
  s := new(Streams)
  s.items = make(map[string]chan string)

  s.Put(name, stream)

  return *s
}

func NewStreams2(name1 string, stream1 chan string, name2 string, stream2 chan string) Streams {
  s := new(Streams)
  s.items = make(map[string]chan string)

  s.Put(name1, stream1)
  s.Put(name2, stream2)

  return *s
}

func NewStreams3(
  name1 string, stream1 chan string,
  name2 string, stream2 chan string,
  name3 string, stream3 chan string,
) Streams {
  s := new(Streams)
  s.items = make(map[string]chan string)

  s.Put(name1, stream1)
  s.Put(name2, stream2)
  s.Put(name3, stream3)

  return *s
}

// TODO: Extract

func (s *Streams) Put(name string, stream chan string) *Streams {
  s.items[name] = stream
  return s
}
