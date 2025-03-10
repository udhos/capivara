#!/bin/bash

build() {
    local branch="$1"
    local bin="$2"
    git checkout $branch
    go build -o ~/go/bin/$bin ./capivara
    $bin -version
    git checkout -
}

build main capivara
build 3fr capivara-3fr
build qs capivara-qs
