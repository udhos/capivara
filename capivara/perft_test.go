package main

import "testing"

func TestPerft(t *testing.T) {

	game := newGame()
	game.loadFromString(perftBoard)
	b := game.history[len(game.history)-1]

	for d, nodes := range testPerftTable {

		if d > 5 {
			break
		}

		buf := defaultBoardPool
		buf.reset()

		n, _ := perft(b, d, buf)

		if n != nodes {
			t.Errorf("perft depth %d: got %d nodes, expected %d", d, n, nodes)
		}
	}
}

const perftBoard = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |.R|.N|.B|.Q|.K|.B|.N|.R|  8
   -------------------------
7  |.p|.p|.p|.p|.p|.p|.p|.p|  7
   -------------------------
6  |  |  |  |  |  |  |  |  |  6
   -------------------------
5  |  |  |  |  |  |  |  |  |  5
   -------------------------
4  |  |  |  |  |  |  |  |  |  4
   -------------------------
3  |  |  |  |  |  |  |  |  |  3
   -------------------------
2  |*p|*p|*p|*p|*p|*p|*p|*p|  2
   -------------------------
1  |*R|*N|*B|*Q|*K|*B|*N|*R|  1
   -------------------------
    a  b  c  d  e  f  g  h
`

func TestPerft2(t *testing.T) {

	game := newGame()
	game.loadFromString(perftBoard2)
	b := game.history[len(game.history)-1]

	d := 2

	buf := defaultBoardPool
	buf.reset()

	n, _ := perft(b, d+1, buf)

	expected := int64(33949)
	if n != expected {
		t.Errorf("perft depth %d: got %d nodes, expected %d", d, n, expected)
	}
}

const perftBoard2 = `
    a  b  c  d  e  f  g  h
   -------------------------
8  |.R|.N|.B|.Q|.K|  |  |.R|  8
   -------------------------
7  |.p|.p|.p|.p|  |.p|.p|.p|  7
   -------------------------
6  |  |  |  |  |  |.N|  |  |  6
   -------------------------
5  |  |  |.B|  |.p|  |  |  |  5
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
