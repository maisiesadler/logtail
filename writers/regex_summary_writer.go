package writers

import (
  "encoding/json"
  "io"
  "regexp"
  "strconv"
)

type RegexSummaryWriter struct {
  w io.Writer
  Summary *Summary
}

func NewRegexSummaryWriter(r *regexp.Regexp, w io.Writer) *RegexSummaryWriter {
  summary := &Summary{linecount: make(map[string]int), r: r}
  return &RegexSummaryWriter{w, summary}
}

func (w *RegexSummaryWriter) Write(data []byte) (n int, err error) {
  if added := w.Summary.Add(string(data)); added {
//     w.w.Write([]byte(strings.Join(w.Summary.ReadAsString(), "\n")))
    w.w.Write(w.Summary.ReadAsJson())
  }
  return 0, nil
}

func (s *Summary) ReadAsJson() []byte {
  lac := s.Read()
  b, _ := json.Marshal(lac)
  return b
}

func (s *Summary) ReadAsString() []string {
  lac := []string{}
  s.RLock()
  for k, v := range s.linecount {
    lac = append(lac, k + " (" + strconv.Itoa(v) + ")")
  }
  s.RUnlock()
  return lac
}
