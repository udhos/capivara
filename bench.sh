#!/bin/bash

(cd ./capivara && go test -benchtime=100x -run=Benchmark -bench=.)
