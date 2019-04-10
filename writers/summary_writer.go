package writers

import (
  "sync"
  "regexp"
)

type SummaryWriter struct {
  updates chan bool
  Summary *Summary
}

func NewSummaryWriter(r *regexp.Regexp) *SummaryWriter {
  summary := &Summary{linecount: make(map[string]int), r: r}
  return &SummaryWriter{make(chan bool), summary}
}

func (w *SummaryWriter) Write(data []byte) (n int, err error) {
  if added := w.Summary.Add(string(data)); added {
    w.updates <- true
  }
  return 0, nil
}

func (w *SummaryWriter) Updates() <-chan bool {
  return w.updates
}

type Summary struct {
  sync.RWMutex
  r *regexp.Regexp
  linecount map[string]int
}

type LineAndCount struct {
  Line string `json:"line"`
  Count int `json:"count"`
}

func (s *Summary) Add(l string) bool {
  if matched := s.r.MatchString(l); matched {
    s.Lock()
    s.linecount[l]++
    s.Unlock()
    return true
  }
  return false
}

func (s *Summary) Read() []LineAndCount {
  lac := []LineAndCount{}
  s.RLock()
  for k, v := range s.linecount {
    lac = append(lac, LineAndCount { k, v })
  }
  s.RUnlock()
  return lac
}

