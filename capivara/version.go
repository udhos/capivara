package main

import (
	"fmt"
	"runtime"
)

const (
	version = "0.8.0"

	// uci: uci protocol
	// ab: alpha-beta search
	// id: iterative deepening
	// pst: piece-square tables
	// z: zobrist hashing
	// 3fr: 3-fold repetition
	// qs: quiescence search
	// pvs: principal variation search
	features = "uci ab id pst"
)

func fullVersion() string {
	return fmt.Sprintf("%s %s %s %s GOMAXPROCS=%d", shortVersion(), runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
}

func shortVersion() string {
	return fmt.Sprintf("%s(%s)", version, features)
}

func showFullVersion() {
	fmt.Printf("capivara version %s\n", fullVersion())
}
