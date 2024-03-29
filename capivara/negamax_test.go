package main

import "testing"

func TestB1(t *testing.T) {
	game := newGame()
	game.loadFromString(b1)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 2, false)
	if m.String() != "d3e4" {
		t.Errorf("score: %v move: %s (expected: move: d3e4)", score, m)
	}
}

func TestB2(t *testing.T) {
	game := newGame()
	game.loadFromString(b2)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 2, false)
	if m.String() != "d4e5" {
		t.Errorf("score: %v move: %s (expected: move: d4e5)", score, m)
	}
}

func TestB3(t *testing.T) {
	game := newGame()
	game.loadFromString(b3)
	last := len(game.history) - 1
	b := game.history[last]
	b.disableCastling()

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, comment := rootNegamax(&nega, b, 2, false)
	if score != -1000.0 || comment != "checkmated" {
		t.Errorf("score: %v move: %s (expected: score=-1000.0 move: checkmated)", score, m)
	}
}

func TestB4(t *testing.T) {
	game := newGame()
	game.loadFromString(b4)
	last := len(game.history) - 1
	b := game.history[last]
	b.turn = colorBlack

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 2, false)
	if m.String() != "e6d5" {
		t.Errorf("score: %v move: %s (expected: move: e6d5)", score, m)
	}
}

func TestB4Depth6(t *testing.T) {
	game := newGame()
	game.loadFromString(b4)
	last := len(game.history) - 1
	b := game.history[last]
	b.turn = colorBlack

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 6, false)
	if m.String() != "e6d5" {
		t.Errorf("score: %v move: %s (expected: move: e6d5)", score, m)
	}
}

func TestB5(t *testing.T) {
	game := newGame()
	game.loadFromString(b5)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 2, false)
	if m.String() != "f6f7" {
		t.Errorf("score: %v move: %s (expected: checkmate f6f7)", score, m)
	}
}

func TestB6(t *testing.T) {
	game := newGame()
	game.loadFromString(b6)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, _, _ := rootNegamax(&nega, b, 2, false)
	if score != 1000 {
		t.Errorf("score: %v (expected: score=1000.0)", score)
	}
}

func TestB7(t *testing.T) {
	game := newGame()
	game.loadFromString(b7)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 4, false)
	if m.String() != "g6g7" {
		t.Errorf("score: %v move: %s (expected: move g6g7)", score, m)
	}
}

func TestB8(t *testing.T) {
	game := newGame()
	game.loadFromString(b8)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 2, false)
	if m.String() != "d2d3" {
		t.Errorf("score: %v move: %s (expected: move d2d3)", score, m)
	}
}

func TestB9(t *testing.T) {
	game := newGame()
	game.loadFromString(b9)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 4, false)
	if m.String() != "g7g8q" {
		t.Errorf("score: %v move: %s (expected: move g7g8q)", score, m)
	}
}

func TestB10(t *testing.T) {
	game := newGame()
	game.loadFromString(b10)
	last := len(game.history) - 1
	b := game.history[last]
	b.turn = colorBlack

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 4, false)
	if m.String() != "e8d7" {
		t.Errorf("score: %v move: %s (expected: move e8d7)", score, m)
	}
}

func TestB11(t *testing.T) {
	game := newGame()
	game.loadFromString(b11)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 4, false)
	if m.String() != "e1g1" {
		t.Errorf("score: %v move: %s (expected: move e1g1)", score, m)
	}
}

func TestB12(t *testing.T) {
	game := newGame()
	game.loadFromString(b12)
	last := len(game.history) - 1
	b := game.history[last]

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children}

	score, m, _ := rootNegamax(&nega, b, 4, false)
	if m.String() != "h1h6" {
		t.Errorf("score: %v move: %s (expected: move h1h6)", score, m)
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
6  |  |  |  |  |  |.p|  |  |  6
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

const b9 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |  |  |  8
   -------------------------
7  |  |  |  |  |.p|  |*p|  |  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
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

const b10 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |  |  |  |  |.K|  |*Q|  |  8
   -------------------------
7  |  |  |  |  |.p|  |  |  |  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
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

const b11 = `
   a  b  c  d  e  f  g  h
   -------------------------
8  |  |.N|.B|  |.R|  |  |  |  8
   -------------------------
7  |.p|.p|.p|.p|.p|  |  |  |  7
   -------------------------
6  |  |  |  |  |  |  |  |.R|  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |*Q|  |  |  |  |  |  |  3
   -------------------------
2  |*p|*p|*p|*p|  |  |  |  |  2
   -------------------------
1  |  |.K|  |  |*K|  |  |*R|  1
   -------------------------
    a  b  c  d  e  f  g  h
`

const b12 = `
   a  b  c  d  e  f  g  h
   -------------------------
8  |  |.N|.B|  |.R|  |  |  |  8
   -------------------------
7  |.p|.p|.p|.p|.p|  |  |  |  7
   -------------------------
6  |  |  |  |  |  |.R|  |.R|  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |*Q|  |  |  |  |  |  |  3
   -------------------------
2  |*p|*p|*p|*p|  |  |  |  |  2
   -------------------------
1  |  |.K|  |  |*K|  |  |*R|  1
   -------------------------
    a  b  c  d  e  f  g  h
`
