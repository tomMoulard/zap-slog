# zap: slog handler


A [slog](https://pkg.go.dev/log/slog) Core for [zap](https://github.com/uber-go/zap) library.

## Install

```bash
go get -u github.com/tommoulard/zap-slog
```

## Usage

GoDoc: [https://pkg.go.dev/github.com/tommoulard/zap-slog](https://pkg.go.dev/github.com/tommoulard/zap-slog)

## Example

```go
package main

import (
	"log/slog"

	zapslog "github.com/tommoulard/zap-slog"
	"go.uber.org/zap"
)

func main() {
	slogger := slog.Default()
	logger, _ := zap.NewProduction(zapslog.WrapCore(slogger))
	logger = logger.Named("example")
	logger.Info("hello world")
}
```

## Benchmarks

Here are a few benchmarks comparing the performance of the `zap-slog` wrapper
against the `zap` library.

Both benchmarks are run with a logger that discards all logs.

```bash
$ go test ./... -bench=.
goos: linux
goarch: amd64
pkg: github.com/tommoulard/zap-slog
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkWrapCore-12    	   677972	          1570 ns/op
BenchmarkZap-12         	217827159	         5.336 ns/op
PASS
ok  	github.com/tommoulard/zap-slog	2.810s
```

