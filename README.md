# api-server
This will be used to host API requests from agent and manage devices.

## Building

### Build with go

- you need go `v1.5` or later.
- set $GOPATH accordingly.

```console
$ go build -o apiserver main.go
```

## Running

- you need mysql-server installed.

```console
$ ./apiserver
```
