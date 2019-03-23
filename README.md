# LogTail

Tails log files.

`go get github.com/maisiesadler/logtail`

Tail and print to console
```
w := &logtail.FmtWriter{}
logtail.Run(file, w)
```

Tail and print to console if matches regex
```
r := regexp.MustCompile("myreg")
f := &logtail.RegexFilter{r}
fw := &logtail.FilteredWriter{&logtail.FmtWriter{}, f}
logtail.Run(file, fw)
```
