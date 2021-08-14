package main

import (
	"strings"
	"testing"
)

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

type perftFENTest struct {
	name          string
	fen           string
	expectedNodes []int64
}

// https://www.chessprogramming.org/Perft_Results
var perftFENTestTable = []perftFENTest{
	{"Initial Position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", testPerftTable},
	{"Position 2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -", []int64{0, 48, 2039, 97862, 4085603, 193690690, 8031647685}},
}

func TestPerftFEN(t *testing.T) {
	testPerftFENDepth(t, 5) // deeper will take long
}

func testPerftFENDepth(t *testing.T, maxDepth int) {

	for _, data := range perftFENTestTable {

		game := newGame()
		fenTokens := strings.Fields(data.fen)
		game.loadFromFen(fenTokens)
		b := game.history[len(game.history)-1]

		buf := defaultBoardPool

		var depth int

		for _, depthNodes := range data.expectedNodes {

			if maxDepth > 0 {
				if depth > maxDepth {
					break
				}
			}

			buf.reset()

			n, _ := perft(b, depth, buf)

			if n != depthNodes {
				t.Errorf("%s: perft maxDepth=%d depth %d: got %d nodes, expected %d", data.name, maxDepth, depth, n, depthNodes)
			} else {
				t.Logf("%s: perft maxDepth=%d depth %d: got %d nodes, expected %d: OK", data.name, maxDepth, depth, n, depthNodes)
			}

			depth++
		}
	}
}
