package logtail

import (
  "log"
  "os"
  "bufio"
  "fmt"
  "sync"
  "strconv"

  "github.com/fsnotify/fsnotify"
)

func Run(file string) {
  notify := make(chan bool)
  done := make(chan bool)
  var wg sync.WaitGroup

  wg.Add(1)

  r := reader(file)
  readUntilEof(r)
  go WatchFile(file, notify, done)
  go func() {
    defer wg.Done()
    for _ = range notify {
      readUntilEof(r)
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

func readUntilEof(r *bufio.Reader) {
  _, err := r.Peek(1)
  for err == nil {
    token, isPrefix, lineerr := r.ReadLine()
    if lineerr != nil {
        panic(lineerr)
    }
    fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
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
