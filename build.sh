#!/bin/bash

go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/tools/cmd/deadcode@latest
go install github.com/mgechev/revive@latest

gofmt -s -w .

revive ./...

go mod tidy

govulncheck ./...

deadcode ./capivara

#export CGO_ENABLED=1

echo "***"
echo "*** TestPerftFEN is slow (it takes about 15 seconds)"
echo "*** but -race is disabled since it became painfully slow in go1.24.1"
echo "***"
#go test -race ./... ;# -race became slow for perft in go1.24.1
go test ./...

export CGO_ENABLED=0

go install ./...
