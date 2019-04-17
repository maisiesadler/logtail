package writers

import (
  "fmt"
  "net/http"
  "strconv"
  "sync"

  "github.com/gorilla/websocket"
)

type WebSocketWriter struct {
  port int
  endpoint string
  subscribers []chan string
  latest []string
  mux *sync.RWMutex
}

func NewWebSocketWriter(port int, endpoint string) *WebSocketWriter {
  updates := []chan string{}
  return &WebSocketWriter{port, endpoint, updates, []string{}, &sync.RWMutex{}}
}

func (w *WebSocketWriter) Write(data []byte) (n int, err error) {
  s := string(data)
  fmt.Println("sending " + s, len(w.subscribers))
  w.mux.Lock()
  w.latest = append(w.latest, s)
  w.mux.Unlock()
  for _, subs := range w.subscribers {
    subs <- s
  }
  return 0, nil
}

func (w *WebSocketWriter) Start() {
  http.HandleFunc("/" + w.endpoint, func(res http.ResponseWriter, req *http.Request) {
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if error != nil {
      http.NotFound(res, req)
      return
    }
    go subscribe(w, conn)
  })
  http.ListenAndServe(":" + strconv.Itoa(w.port), nil)
  fmt.Println("Listening")
}

func subscribe(w *WebSocketWriter, conn *websocket.Conn) {
  defer conn.Close()
  writequeue := make(chan string)
  defer close(writequeue)
  go func() {
    l := 200
    lenl := len(w.latest)
    if lenl < 200 {
      l = lenl
    }
    for _, message := range w.latest[lenl - l: lenl] {
      fmt.Println("latest" + message)
      writequeue <- message
    }
  }()
  w.subscribers = append(w.subscribers, writequeue)
  for u := range writequeue {
    fmt.Println(u)
    conn.WriteMessage(websocket.TextMessage, []byte(u))
  }
}
