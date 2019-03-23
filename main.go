package logtail

import (
  "log"
  "os"
  "bufio"
  "fmt"
  "sync"
  "strconv"
  "io"

  "github.com/fsnotify/fsnotify"
)

type FmtWriter struct {
  io.Writer
}

func (f *FmtWriter) Write(data []byte) (n int, err error) {
  fmt.Println(string(data))
  return 0, nil
}

func Run(file string, w io.Writer) {
  notify := make(chan bool)
  done := make(chan bool)
  var wg sync.WaitGroup

  wg.Add(1)

  r := reader(file)
  readUntilEof(r, w)
  go WatchFile(file, notify, done)
  go func() {
    defer wg.Done()
    for _ = range notify {
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
    // todo: isPrefix?
//     fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
    _, err = r.Peek(1)
  }
  if err != nil {
    fmt.Println(err)
  }
}

func WatchFile(file string, notify chan<- bool, done <-chan bool) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }
  defer watcher.Close()

  go func() {
    for {
      select {
      case event, ok := <-watcher.Events:
        if !ok {
          return
        }
        log.Println("event:", event)
        if event.Op&fsnotify.Write == fsnotify.Write {
          notify <- true
          log.Println("modified file:", event.Name)
        }
      case err, ok := <-watcher.Errors:
        if !ok {
          return
        }
        log.Println("error:", err)
      }
    }
  }()

  err = watcher.Add(file)
  if err != nil {
    log.Fatal(err)
  }
  <-done
}

func AppendFile(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
      panic(err)
  }
  defer f.Close()
  for i := 0; i < 1000; i++ {
    text := "testing" + strconv.Itoa(i) + "\n"
    if _, err = f.WriteString(text); err != nil {
      panic(err)
    }
  }
}
