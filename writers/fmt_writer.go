package writers

import (
  "fmt"
  "io"
)

type FmtWriter struct {
  io.Writer
}

func (f *FmtWriter) Write(data []byte) (n int, err error) {
  fmt.Println(string(data))
  return 0, nil
}

