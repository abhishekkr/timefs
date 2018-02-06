## timefs

> a file system not for regular files but time series data, so it can juice out real efficiency

---

### Usage

* running dummy client

```
## to prepare env
source go-tasks
./go-tasks deps

## to start timefs server
go run server/tfserver.go

## to push in 1000 timedots
go run client/tfclient.go --server="127.0.0.1:7999" dummy create

## to read 1000 timedots
go run client/tfclient.go --server="127.0.0.1:7999" dummy read
```

---

### Performance Metrics

> server and cli running on same node, speed mainly depends on IOPS in current strategy

* on `xfs` cloud 10GB non-SSD

```
writes: 3240 timedots/sec
reads:  21600 timedots/sec
```

* on `ext4` localhost SSD

```
writes: 21600 timedots/sec
reads:  51840 timedots/sec
```

---
---

