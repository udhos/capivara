package main

import "testing"

func TestB1(t *testing.T) {
	game := newGame()
	game.loadFromString(b1)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if score != 3.0 || move != "d2e3" {
		t.Errorf("score: %v move: %s (expected: score=3.0 move: d2e3)", score, move)
	}
}

func TestB2(t *testing.T) {
	game := newGame()
	game.loadFromString(b2)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if score != -2.0 || move != "d4e5" {
		t.Errorf("score: %v move: %s (expected: score=-2.0 move: d4e5)", score, move)
	}
}

func TestB3(t *testing.T) {
	game := newGame()
	game.loadFromString(b3)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if score != -1000.0 || move != "checkmated" {
		t.Errorf("score: %v move: %s (expected: score=-1000.0 move: checkmated)", score, move)
	}
}

func TestB4(t *testing.T) {
	game := newGame()
	game.loadFromString(b4)
	last := len(game.history) - 1
	b := game.history[last]
	b.turn = colorInverse(b.turn)
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if score != 3.0 || move != "e5d4" {
		t.Errorf("score: %v move: %s (expected: score=3.0 move: e5d4)", score, move)
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

const b2 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |*p|  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |.p|  |  |  |  5
   -------------------------
4  |  |  |  |*p|  |*K|  |  |  4
   -------------------------
3  |  |.R|  |  |  |  |  |  |  3
   -------------------------
2  |*p|  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |  |  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b3 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |*p|  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |.p|  |  |  |  5
   -------------------------
4  |  |  |  |  |  |*K|  |  |  4
   -------------------------
3  |  |.R|  |  |  |  |  |  |  3
   -------------------------
2  |*p|  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |  |  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b4 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |.p|  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |.p|  |  |  |  5
   -------------------------
4  |  |  |  |*p|  |  |  |  |  4
   -------------------------
3  |  |  |  |  |  |  |  |  |  3
   -------------------------
2  |.p|  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |*K|  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`
