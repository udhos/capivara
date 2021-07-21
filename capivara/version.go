package main

import (
	"fmt"
	"runtime"
)

const version = "0.8.0"

func fullVersion() string {
	return fmt.Sprintf("%s %s %s %s GOMAXPROCS=%d", version, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
}
