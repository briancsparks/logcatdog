package cdr0

type Signals struct {
  items map[string]chan struct{}
}

func NewSignals() Signals {
  s := new(Signals)
  s.items = make(map[string]chan struct{})
  return *s
}

func (s *Signals) SafeSignal(name string) chan struct{} {
  result := s.items[name]
  if result != nil {
    return result
  }

  return make(chan struct{})
}

// TODO: Extract

func (s *Signals) Put0(name string) chan struct{} {
  signal := make(chan struct{})
  s.items[name] = signal
  return signal
}


