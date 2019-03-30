# LogTail

Tails log files.

`go get github.com/maisiesadler/logtail`
`go get github.com/maisiesadler/logtail/writers`

Tail and print to console
```
w := &writers.FmtWriter{}
logtail.Run(file, w)
```

Tail and print to console if matches regex
```
r := regexp.MustCompile("myreg")
f := &writers.RegexFilter{r}
fw := &writers.FilteredWriter{&logtail.FmtWriter{}, f}
logtail.Run(file, fw)
```

Tail and write to websocket connection
```
w := writers.NewWebSocketWriter()
go logtail.Run(file, w)
w.Start()
```

Tail, filter using regex, then write to websocket connection
```
w := writers.NewWebSocketWriter()
f := &writers.RegexFilter{regexp.MustCompile("myreg")}
fw := &writers.FilteredWriter{w, f}
go logtail.Run(file, fw)
w.Start()
```
