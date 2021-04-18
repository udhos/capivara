package main

import "testing"

var move string

func BenchmarkCastling(b *testing.B) {
	game := newGame()
	game.loadFromString(castling)
	brd := game.history[len(game.history)-1]
	ab := alphaBetaState{}
	var mv string
	for n := 0; n < b.N; n++ {
		_, m, _ := rootAlphaBeta(&ab, brd, 2, nil, false)
		mv = m // record call result to prevent compiler from eliminating function call
	}
	move = mv // record bench result to prevent the compiler from eliminating the test
}

func BenchmarkCastlingAddChildren(b *testing.B) {
	game := newGame()
	game.loadFromString(castling)
	brd := game.history[len(game.history)-1]
	ab := alphaBetaState{}
	var mv string
	for n := 0; n < b.N; n++ {
		_, m, _ := rootAlphaBeta(&ab, brd, 2, nil, true)
		mv = m // record call result to prevent compiler from eliminating function call
	}
	move = mv // record bench result to prevent the compiler from eliminating the test
}

const castling = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |.R|  |.B|.Q|.K|.B|  |.R|  8
   -------------------------
7  |.p|.p|.p|.p|  |.p|.p|.p|  7
   -------------------------
6  |  |  |.N|  |  |.N|  |  |  6
   -------------------------
5  |  |  |  |  |.p|  |  |  |  5
   -------------------------
4  |  |  |*B|  |*p|  |  |  |  4
   -------------------------
3  |  |  |  |  |  |*N|  |  |  3
   -------------------------
2  |*p|*p|*p|*p|  |*p|*p|*p|  2
   -------------------------
1  |*R|*N|*B|*Q|*K|  |  |*R|  1
   -------------------------
    a  b  c  d  e  f  g  h
`
