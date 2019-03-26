package writers

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

