#!/bin/bash

(cd ./capivara && go test -benchtime=10x -run=Benchmark -bench=.)
