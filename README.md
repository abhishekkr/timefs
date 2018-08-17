## timefs

> a file system not for regular files but time series data, so it can juice out real efficiency

here, we'll see a term `timedot` in use, a `timedot` is every individual value persisted at any individual timestamp

`
---

### HowTo Use

> can be used with single server, or multiple servers fronted by a request splitter
>
> current splitter is really basic and under improvements

* to activate local go path just for your editor or direct running go tools

`. go-tasks`

* to prepare env, source go-tasks and then

`./go-tasks deps`

* start backends at ':8001' and ':8002'

```
## in one session
TIMEFS_DIR_ROOT=/tmp/timefs/backend0 TIMEFS_PORT=":8001" go run server/tfserver.go

## in other session
TIMEFS_DIR_ROOT=/tmp/timefs/backend1 TIMEFS_PORT=":8002" go run server/tfserver.go
```


* start splitter to manage all backends ':8001,:8002'

```
TIMEFS_CLIENTBYCHANNEL_COUNT=10 TIMEFS_PROXY_PORT=":7999" TIMEFS_BACKENDS="127.0.0.1:8001,127.0.0.1:8002" go run splitter/tfsplitter.go
```


* start client to use backends via splitter

```
go run client/tfclient.go --server="127.0.0.1:7999" dummy create

go run client/tfclient.go --server="127.0.0.1:7999" dummy read
```



[here you can check detailed overview on usage options](./docs/usage.md)

---

### Performance Metrics

##### with current splitter

requests for 48000 timedots, rough estimates

* one backend with `writes:	0m2.47s` and `reads: 0m0.36s`

* two backends with `writes: 0m2.46s` and `reads: 0m0.18s`

* three backends with `writes: 0m2.45s` and `reads: 0m0.24s`

---

