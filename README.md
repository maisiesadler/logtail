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

Example Client App
```
<script>

    function reconnectAfterTimeout() {
        //     const backoff = 1 * 1000 + 500;
        //     setTimeout(() => this.reconnect(), backoff);
        //     if (this.backoff < 30) {
        //         this.backoff++;
        //     }
    }

    function createSocket(url, onmessage, onerror, onclose, onopen) {
        const queue = [];
        const ws = new WebSocket(url);
        ws.onmessage = msg => typeof onmessage === 'function' && onmessage(msg)
        ws.onerror = err => typeof onerror === 'function' && onerror(err)
        ws.onclose = err => typeof onclose === 'function' && onclose(err)
        ws.onopen = () => {
            console.log('opened, queue=', queue)
            while (queue.length > 0) {
                let a = queue.shift();
                ws.send(a);
                console.log('sent queued message', a)
            }
            typeof onopen === 'function' && onopen(err)
        }

        const send = msg => {
            console.log(ws);
            if (ws.readyState === WebSocket.OPEN) {
                console.log('sending', msg);
                // ws.send(JSON.stringify(data));
                ws.send(msg);
            } else if (ws.readyState === WebSocket.CLOSED) {
                console.log('not sending, closed, queuing', msg);
                queue.push(msg);
            } else {
                console.log('not ready, queuing');
                queue.push(msg);
            }
        };
        return { send, close: ws.close };
    }

    function appendMessage(message) {
        const newel = document.createElement('div')
        newel.append(message)
        newel.className = 'message'
        document.querySelector('#messages').append(newel)
        const messageContainer = document.querySelector('#message-container')
        messageContainer.scrollTop = messageContainer.scrollHeight - messageContainer.clientHeight;
    }

    function replaceMessage(message) {
        const list = JSON.parse(message)
        document.querySelector('#messages').innerHTML = ''
        list.forEach(element => {
            const newel = document.createElement('div')
            newel.append(`${element.line} - ${element.count}`)
            newel.className = 'message'
            document.querySelector('#messages').append(newel)
        });

        const messageContainer = document.querySelector('#message-container')
        messageContainer.scrollTop = messageContainer.scrollHeight - messageContainer.clientHeight;
    }

    function reconnect() {
        console.log('trying again');
        const wsUrl = 'ws://localhost:8080/';
        const onmessage = msg => {
            replaceMessage(msg.data);
        };
        const onerror = err => reconnectAfterTimeout();
        const onclose = () => reconnectAfterTimeout();
        const socket = createSocket(wsUrl + 'echo', onmessage, onerror, onclose);
    }

    reconnect()
</script>

<style>
    html,
    body {
        height: 100vh;
    }

    #message-container {
        height: calc(100vh - 70px);
        position: relative;
        overflow: scroll;
    }

    #messages {
        max-height: 100%;
        bottom: 0;
    }
</style>

<h1>Logs</h1>
<div id="message-container">
    <div id="messages">

    </div>
</div>
```
