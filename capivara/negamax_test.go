package main

import "testing"

func TestB1(t *testing.T) {
	game := newGame()
	game.loadFromString(b1)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if move != "d3e4" {
		t.Errorf("score: %v move: %s (expected: move: d3e4)", score, move)
	}
}

func TestB2(t *testing.T) {
	game := newGame()
	game.loadFromString(b2)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if move != "d4e5" {
		t.Errorf("score: %v move: %s (expected: move: d4e5)", score, move)
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
	b.turn = colorBlack
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if move != "e6d5" {
		t.Errorf("score: %v move: %s (expected: move: e6d5)", score, move)
	}
}

func TestB4Depth6(t *testing.T) {
	game := newGame()
	game.loadFromString(b4)
	last := len(game.history) - 1
	b := game.history[last]
	b.turn = colorBlack
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 6, []string{})
	if move != "e6d5" {
		t.Errorf("score: %v move: %s (expected: move: e6d5)", score, move)
	}
}

func TestB5(t *testing.T) {
	game := newGame()
	game.loadFromString(b5)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 2, []string{})
	if move != "f6f7" {
		t.Errorf("score: %v move: %s (expected: checkmate f6f7)", score, move)
	}
}

func TestB6(t *testing.T) {
	game := newGame()
	game.loadFromString(b6)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, _, _ := rootNegamax(&nega, b, 2, []string{})
	if score != 1000 {
		t.Errorf("score: %v (expected: score=1000.0)", score)
	}
}

func TestB7(t *testing.T) {
	game := newGame()
	game.loadFromString(b7)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 4, []string{})
	if move != "g6g7" {
		t.Errorf("score: %v (expected: move g6g7)", score)
	}
}

func TestB8(t *testing.T) {
	game := newGame()
	game.loadFromString(b8)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, _ := rootNegamax(&nega, b, 4, []string{})
	if move != "d2d3" {
		t.Errorf("score: %v (expected: move d2d3)", score)
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
4  |  |  |  |  |.p|  |  |  |  4
   -------------------------
3  |  |  |  |*p|  |  |  |  |  3
   -------------------------
2  |  |  |  |  |  |  |  |  |  2
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
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |.R|  |  |  |.p|  |  |  3
   -------------------------
2  |*p|  |  |  |  |  |.p|*N|  2
   -------------------------
1  |  |  |  |  |  |  |*B|*K|  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b4 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |  |  7
   -------------------------
6  |  |  |  |  |.p|  |  |  |  6
   -------------------------
5  |  |  |  |*p|  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |  |  |  |  |  |  |  |  3
   -------------------------
2  |  |  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |*K|  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b5 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |.R|.p|  7
   -------------------------
6  |  |  |  |  |  |*p|  |  |  6
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

const b6 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |*p|  |  |  7
   -------------------------
6  |  |  |  |.p|  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |*p|  |  |  |  |  |  4
   -------------------------
3  |  |  |  |  |  |  |  |  |  3
   -------------------------
2  |  |  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |*K|  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b7 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |  |  |  |  |  7
   -------------------------
6  |  |  |  |  |  |  |*p|  |  6
   -------------------------
5  |  |  |  |  |  |  |  |*p|  5
   -------------------------
4  |  |  |.p|  |  |  |  |  |  4
   -------------------------
3  |  |.p|  |  |  |  |  |  |  3
   -------------------------
2  |  |  |  |  |  |  |  |  |  2
   -------------------------
1  |  |  |  |  |*K|  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b8 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |.K|  |  |  |  |  8
   -------------------------
7  |  |  |  |  |.p|*B|  |  |  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |  |  |.p|  |*N|  |  |  3
   -------------------------
2  |  |  |  |*R|  |  |  |*K|  2
   -------------------------
1  |  |  |  |  |  |  |  |  |  1
   -------------------------
    a  b  c  d  e  f  g  h

`
