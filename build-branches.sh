#!/bin/bash

build() {
    local branch="$1"
    local bin="$2"
    git checkout $branch
    git pull
    go build -o ~/go/bin/$bin ./capivara
    $bin -version
    git checkout -
}

build main capivara
build 3fr capivara-3fr
build 3fr-b capivara-3fr-b
build 3fr-c capivara-3fr-c
build 3fr-d capivara-3fr-d
