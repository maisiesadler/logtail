package logtail

import (
  "log"
  "os"
  "bufio"
  "fmt"
  "sync"
  "io"
  "regexp"
  "time"
)

func Run(file string, w io.Writer) {
  var wg sync.WaitGroup

  wg.Add(1)
  ticker := time.NewTicker(500 * time.Millisecond)

  r := reader(file)
  readUntilEof(r, w)
  go func() {
    defer wg.Done()
    for _ = range ticker.C {
      readUntilEof(r, w)
    }
  }()
  wg.Wait()
}

type FmtWriter struct {
  io.Writer
}

func (f *FmtWriter) Write(data []byte) (n int, err error) {
  fmt.Println(string(data))
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

type ChannelWriter struct {
  updates chan string
}

func NewChannelWriter() *ChannelWriter {
  return &ChannelWriter{ make(chan string) }
}

func (w *ChannelWriter) Write(data []byte) (n int, err error) {
  w.updates <- string(data)
  return 0, nil
}

func (w *ChannelWriter) Updates() <-chan string {
  return w.updates
}

func reader(filename string) *bufio.Reader {
  f, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  r := bufio.NewReader(f)
  return r
}

func readUntilEof(r *bufio.Reader, w io.Writer) {
  _, err := r.Peek(1)
  for err == nil {
    token, _, lineerr := r.ReadLine()
    if lineerr != nil {
        panic(lineerr)
    }
    w.Write(token)
    _, err = r.Peek(1)
  }
  if err != nil && err != io.EOF {
    log.Println(err)
  }
}

