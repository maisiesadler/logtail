package logtail

import (
  "log"
  "os"
  "bufio"
  "io"
  "time"
)

func Run(file string, w io.Writer) chan<- bool {
  ticker := time.NewTicker(500 * time.Millisecond)

  r := reader(file)
  readUntilEof(r, w)
  go func() {
    for _ = range ticker.C {
      readUntilEof(r, w)
    }
  }()

  done := make(chan bool)
  go func() {
    <-done
    ticker.Stop()
  }()
  return done
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

