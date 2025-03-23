package main

import "testing"

var testMove move

func BenchmarkCastling(b *testing.B) {
	game := newGame()
	game.loadFromString(castling)
	brd := &game.history[len(game.history)-1]

	children := defaultBoardPool
	ab := alphaBetaState{children: children}

	var mv move
	for n := 0; n < b.N; n++ {
		children.reset()
		_, m, _ := rootAlphaBeta(&ab, brd, 2, false)
		mv = m // record call result to prevent compiler from eliminating function call
	}
	testMove = mv // record bench result to prevent the compiler from eliminating the test
}

func BenchmarkCastlingAddChildren(b *testing.B) {
	game := newGame()
	game.loadFromString(castling)
	brd := &game.history[len(game.history)-1]

	children := defaultBoardPool
	ab := alphaBetaState{children: children}

	var mv move
	for n := 0; n < b.N; n++ {
		children.reset()
		_, m, _ := rootAlphaBeta(&ab, brd, 2, true)
		mv = m // record call result to prevent compiler from eliminating function call
	}
	testMove = mv // record bench result to prevent the compiler from eliminating the test
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

// go test -run TestRepetition ./capivara
func TestRepetition(t *testing.T) {

	zobristInit()

	game := newGame()
	game.loadFromString(builtinBoard)

	z1 := repetitionPlay(t, &game, "g1f3")
	repetitionPlay(t, &game, "g8f6")
	repetitionPlay(t, &game, "b1c3")
	repetitionPlay(t, &game, "f6g8")
	z2 := repetitionPlay(t, &game, "c3b1")

	if z1 != z2 {
		t.Errorf("z1=%s != z2=%s", z1, z2)
	}
}

func repetitionPlay(t *testing.T, g *gameState, moveStr string) zobristKey {
	b := &g.history[len(g.history)-1]
	t.Logf("%s: before: zobrist: %s reversible: %t repetition: %t",
		moveStr, b.zobristValue, b.reversible, b.isRepetition())
	if err := g.play(moveStr); err != nil {
		t.Fatal(err)
	}
	b = &g.history[len(g.history)-1]
	t.Logf("%s: after : zobrist: %s reversible: %t, repetition: %t",
		moveStr, b.zobristValue, b.reversible, b.isRepetition())
	return b.zobristValue
}
