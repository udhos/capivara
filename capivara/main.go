package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"unicode"
	"unsafe"
)

type gameState struct {
	history     []board
	addChildren bool
	cpuprofile  string
}

func (g *gameState) play(moveStr string) error {
	last := len(g.history) - 1
	b := g.history[last]

	children := defaultBoardPool
	children.reset()
	b.generateChildren(children)

	m, errMove := newMove(moveStr)
	if errMove != nil {
		return errMove
	}

	for _, c := range children.pool {
		if c.lastMove == m {
			// found valid move
			g.history = append(g.history, c)
			return nil
		}
	}
	return fmt.Errorf("not a valid move: %s", moveStr)
}

func newGame() gameState {
	return gameState{history: []board{{}}}
}

func (g gameState) show() {
	b := g.history[len(g.history)-1] // read-only copy
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Println("   -------------------------")
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d  |", row+1)
		for col := 0; col < 8; col++ {
			loc := row*8 + col
			piece := b.square[loc]
			piece.show()
			fmt.Print("|")
		}
		fmt.Printf("  %d\n", row+1)
		fmt.Println("   -------------------------")
	}
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Printf("turn: %s\n", b.turn.name())

	children := defaultBoardPool
	children.reset()

	fmt.Printf("material: %v evaluation: %v\n", b.getMaterialValue(), relativeMaterial(children, b, g.addChildren))
	fmt.Printf("white king=%s material=%d castlingLeft=%v castlingRight=%v\n", locToStr(b.king[0]), b.materialValue[0], b.flags[0]&lostCastlingLeft == 0, b.flags[0]&lostCastlingRight == 0)
	fmt.Printf("black king=%s material=%d castlingLeft=%v castlingRight=%v\n", locToStr(b.king[1]), b.materialValue[1], b.flags[1]&lostCastlingLeft == 0, b.flags[1]&lostCastlingRight == 0)
	g.showFen()
	fmt.Printf("history %d moves: ", len(g.history))
	for _, m := range g.history {
		fmt.Printf("(%s)", m.lastMove)
	}
	fmt.Println()

	children.reset()
	countChildren := b.generateChildren(children)

	fmt.Printf("%d valid moves:", countChildren)
	for _, c := range children.pool {
		fmt.Printf(" %s", c.lastMove)
	}
	fmt.Println()
}

const builtinBoard = `
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

func (g *gameState) loadFromFen(fen []string) {
	b, errFen := fenParse(fen)
	if errFen != nil {
		fmt.Printf("loadFromFen: %v\n", errFen)
		return
	}
	g.history = []board{b} // replace board
}

func (g *gameState) loadFromString(s string) {
	g.loadFromReader(strings.NewReader(s))
}

func (g *gameState) loadFromFile(filename string) {
	fmt.Printf("loadFromFile: %s\n", filename)
	input, errOpen := os.Open(filename)
	if errOpen != nil {
		fmt.Printf("loadFromFile: %s: %v\n", filename, errOpen)
		return
	}
	defer input.Close()
	g.loadFromReader(input)
}

func (g *gameState) loadFromReader(input io.Reader) {

	reader := bufio.NewReader(input)

	var lineCount int
	b := board{} // new board

	for {
		lineCount++
		line, errRead := reader.ReadString('\n')
		switch errRead {
		case io.EOF:
			fmt.Println("loadFromReader: resetting board")
			g.history = []board{b} // replace board
			return
		case nil:
		default:
			fmt.Printf("load error: %v\n", errRead)
			return
		}

		line = strings.TrimSpace(line)

		//fmt.Printf("load: %s line=%d: [%s]\n", filename, lineCount, line)

		row := -1
		col := -1
		var color pieceColor
		for _, c := range line {
			//fmt.Printf("load: %s line=%d char=[%s]\n", filename, lineCount, string(c))
			switch {
			case unicode.IsDigit(c):
				row = int(c) - '0' - 1
			case c == '|':
				col++
			case c == '*':
				color = colorWhite
			case c == '.':
				color = colorBlack
			case c == 'p':
				b.loadPiece(location(row), location(col), piece(color<<3)+whitePawn)
			case c == 'R':
				b.loadPiece(location(row), location(col), piece(color<<3)+whiteRook)
			case c == 'N':
				b.loadPiece(location(row), location(col), piece(color<<3)+whiteKnight)
			case c == 'B':
				b.loadPiece(location(row), location(col), piece(color<<3)+whiteBishop)
			case c == 'Q':
				b.loadPiece(location(row), location(col), piece(color<<3)+whiteQueen)
			case c == 'K':
				b.loadPiece(location(row), location(col), piece(color<<3)+whiteKing)
			}
		}
	}
}

func (b *board) loadPiece(row, col location, p piece) {
	//log.Printf("loadPiece: %c%c %s %s", col+'a', row+'1', p.color().name(), p.kindLetterLow())
	b.addPiece(row, col, p)
}

const version = "0.2.0"

func main() {

	fmt.Printf("capivara version %s runtime %s GOMAXPROCS=%d OS=%s arch=%s\n",
		version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH)

	var positionalTable bool
	var addChildren bool
	var cpuprofile string

	flag.BoolVar(&positionalTable, "positionalTable", positionalTable, "use positional table")
	flag.BoolVar(&addChildren, "addChildren", addChildren, "compute number of children into evalution function")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "save cpuprofile into to file")

	flag.Parse()

	if !positionalTable {
		// clear table
		positionTable = [6][64]int16{
			noWeight, // king
			noWeight, // queen
			noWeight, // rook
			noWeight, // bishop
			noWeight, // knight
			noWeight, // queen
		}
	}

	gameLoop(addChildren, cpuprofile)
}

func gameLoop(addChildren bool, cpuprofile string) {

	game := newGame()
	game.addChildren = addChildren
	game.cpuprofile = cpuprofile
	game.loadFromString(builtinBoard)

	fmt.Printf("board size: %d bytes\n", unsafe.Sizeof(board{}))

	input := bufio.NewReader(os.Stdin)
LOOP:
	for {
		game.show()
		fmt.Print("enter command:")
		text, errInput := input.ReadString('\n')
		switch errInput {
		case io.EOF:
			fmt.Println("input EOF, bye.")
			break LOOP
		case nil:
		default:
			fmt.Printf("input error: %v\n", errInput)
			continue
		}

		tokens := strings.Fields(text)
		if len(tokens) < 1 {
			continue
		}
		cmdPrefix := tokens[0]

		for _, cmd := range tableCmd {
			if strings.HasPrefix(cmd.name, cmdPrefix) {
				cmd.call(tableCmd, &game, tokens)
				continue LOOP
			}
		}

		fmt.Printf("bad command: %s\n", cmdPrefix)
	}
}
