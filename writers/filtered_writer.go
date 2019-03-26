package writers

import (
  "io"
  "regexp"
)

type FilteredWriter struct {
  W io.Writer
  F Filter
}

func (w *FilteredWriter) Write(data []byte) (n int, err error) {
  if ok, d := w.F.Pass(data); ok {
    return w.W.Write(d)
  }
  return 0, nil
}

type Filter interface {
  Pass(data []byte) (bool, []byte)
}

type RegexFilter struct {
  R *regexp.Regexp
}

func (rf *RegexFilter) Pass(data []byte) (bool, []byte) {
  if (rf.R.Match(data)) {
    return true, data
  }
  return false, []byte{}
}

