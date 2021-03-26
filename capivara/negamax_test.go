package main

import "testing"

func TestB1(t *testing.T) {
	game := newGame()
	game.loadFromString(b1)
	last := len(game.history) - 1
	b := game.history[last]
	score, move := rootNegamax(b, 2)
	if score != 3.0 || move != "d2 e3" {
		t.Errorf("score: %v move: %s (expected: score=3.0 move: d2 e3)", score, move)
	}
}

const b1 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |  |  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |  |  |  |.p|  |  |  |  3
   -------------------------
2  |*p|  |  |*p|  |  |*p|  |  2
   -------------------------
1  |  |  |  |  |*K|  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`
