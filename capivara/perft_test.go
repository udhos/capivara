package main

import "testing"

var testPerftTable = []int{0, 20, 400, 8902, 197281, 4865609}

func TestPerft(t *testing.T) {

	game := newGame()
	game.loadFromString(perftBoard)
	b := game.history[len(game.history)-1]

	for d, nodes := range testPerftTable {

		buf := []board(nil)
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
