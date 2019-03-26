package logtail

import (
  "log"
  "os"
  "bufio"
  "sync"
  "io"
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

