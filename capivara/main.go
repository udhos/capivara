package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

type gameState struct {
	history []board
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
	fmt.Printf("material: %v\n", b.getMaterialValue())
	fmt.Printf("white king=%s material=%d castlingLeft=%v castlingRight=%v\n", locToStr(b.king[0]), b.materialValue[0], b.flags[0]&lostCastlingLeft == 0, b.flags[0]&lostCastlingRight == 0)
	fmt.Printf("black king=%s material=%d castlingLeft=%v castlingRight=%v\n", locToStr(b.king[1]), b.materialValue[1], b.flags[1]&lostCastlingLeft == 0, b.flags[1]&lostCastlingRight == 0)
	fmt.Printf("history %d moves: ", len(g.history))
	for _, m := range g.history {
		fmt.Printf("(%s)", m.lastMove)
	}
	fmt.Println()
	children := b.generateChildren([]board{})
	fmt.Printf("%d valid moves:", len(children))
	for _, c := range children {
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
				b.addPiece(location(row), location(col), piece(color<<3)+whitePawn)
			case c == 'R':
				b.addPiece(location(row), location(col), piece(color<<3)+whiteRook)
			case c == 'N':
				b.addPiece(location(row), location(col), piece(color<<3)+whiteKnight)
			case c == 'B':
				b.addPiece(location(row), location(col), piece(color<<3)+whiteBishop)
			case c == 'Q':
				b.addPiece(location(row), location(col), piece(color<<3)+whiteQueen)
			case c == 'K':
				b.addPiece(location(row), location(col), piece(color<<3)+whiteKing)
			}
		}
	}
}

const version = "0.0"

func main() {

	fmt.Printf("capivara version %s runtime %s GOMAXPROCS=%d OS=%s arch=%s\n",
		version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH)

	gameLoop()
}

func gameLoop() {

	game := newGame()
	game.loadFromString(builtinBoard)

	fmt.Printf("board size: %d bytes\n", unsafe.Sizeof(board{}))

LOOP:
	for {
		input := bufio.NewReader(os.Stdin)
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

type command struct {
	name        string
	call        func(cmds []command, game *gameState, tokens []string)
	description string
}

var tableCmd = []command{
	{"ab", cmdAlphaBeta, "alpha-beta search"},
	{"clear", cmdClear, "erase board"},
	{"help", cmdHelp, "show help"},
	{"load", cmdLoad, "load board from file"},
	{"move", cmdMove, "change piece position"},
	{"negamax", cmdNegamax, "negamax search"},
	{"play", cmdPlay, "play move"},
	{"reset", cmdReset, "reset board to initial position"},
	{"switch", cmdSwitch, "switch turn"},
	{"undo", cmdUndo, "undo last played move"},
}

func cmdClear(cmds []command, game *gameState, tokens []string) {
	*game = newGame()
}

func cmdHelp(cmds []command, game *gameState, tokens []string) {
	fmt.Println("available commands:")
	for _, cmd := range cmds {
		fmt.Printf(" %s - %s\n", cmd.name, cmd.description)
	}
}

func cmdLoad(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: load filename\n")
		return
	}
	game.loadFromFile(tokens[1])
}

func cmdMove(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: move fromto\n")
		return
	}
	move := tokens[1]
	if len(move) < 4 || len(move) > 5 {
		fmt.Printf("usage: bad move length=%d: '%s'\n", len(move), move)
		return
	}
	from := strings.ToLower(move[:2])
	to := strings.ToLower(move[2:4])

	fmt.Printf("[%s][%s]\n", from, to)

	// from
	if len(from) != 2 {
		fmt.Printf("bad source format: [%s]\n", from)
		return
	}
	if from[0] < 'a' || from[0] > 'h' {
		fmt.Printf("bad source column letter: [%s]\n", from)
	}
	if from[1] < '1' || from[1] > '8' {
		fmt.Printf("bad source row digit: [%s]\n", from)
	}

	// to
	if len(to) != 2 {
		fmt.Printf("bad target format: [%s]\n", to)
		return
	}
	if to[0] < 'a' || to[0] > 'h' {
		fmt.Printf("bad target column letter: [%s]\n", to)
	}
	if to[1] < '1' || to[1] > '8' {
		fmt.Printf("bad target row digit: [%s]\n", to)
	}

	b := &game.history[len(game.history)-1]                       // will update in-place
	p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board

	if len(move) > 4 {
		// promotion
		promotion := move[4]
		kind := pieceKindFromLetter(rune(promotion))
		if kind != pieceNone {
			p = piece(b.turn<<3) + kind
		}
	}

	b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p) // put piece on board
}

func cmdNegamax(cmds []command, game *gameState, tokens []string) {
	depth := 4
	if len(tokens) > 1 {
		d, errConv := strconv.Atoi(tokens[1])
		if errConv == nil {
			depth = d
		}
	}
	fmt.Printf("negamax depth=%d\n", depth)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, path := rootNegamax(&nega, b, depth, make([]string, 0, 20))
	fmt.Printf("negamax: nodes=%d best score=%v move: %s path: %s\n", nega.nodes, score, move, path)
}

func cmdAlphaBeta(cmds []command, game *gameState, tokens []string) {
	depth := 4
	if len(tokens) > 1 {
		d, errConv := strconv.Atoi(tokens[1])
		if errConv == nil {
			depth = d
		}
	}
	fmt.Printf("alphabeta depth=%d\n", depth)
	last := len(game.history) - 1
	b := game.history[last]
	ab := alphaBetaState{showSearch: true}
	score, move, path := rootAlphaBeta(&ab, b, depth, make([]string, 0, 20))
	fmt.Printf("alphabeta: nodes=%d best score=%v move: %s path: %s\n", ab.nodes, score, move, path)
}

func cmdPlay(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: play fromto\n")
		return
	}
	move := tokens[1]

	// valid move?
	last := len(game.history) - 1
	b := game.history[last]
	if !b.validMove(move) {
		fmt.Printf("not a valid move: %s\n", move)
		return
	}

	for _, c := range b.generateChildren(nil) {
		if c.lastMove == move {
			game.history = append(game.history, c)
			break
		}
	}

	/*
		if len(move) < 4 || len(move) > 5 {
			fmt.Printf("usage: bad move length=%d: '%s'\n", len(move), move)
			return
		}
		from := strings.ToLower(move[:2])
		to := strings.ToLower(move[2:4])

		fmt.Printf("[%s][%s]\n", from, to)

		// from
		if len(from) != 2 {
			fmt.Printf("bad source format: [%s]\n", from)
			return
		}
		if from[0] < 'a' || from[0] > 'h' {
			fmt.Printf("bad source column letter: [%s]\n", from)
		}
		if from[1] < '1' || from[1] > '8' {
			fmt.Printf("bad source row digit: [%s]\n", from)
		}

		// to
		if len(to) != 2 {
			fmt.Printf("bad target format: [%s]\n", to)
			return
		}
		if to[0] < 'a' || to[0] > 'h' {
			fmt.Printf("bad target column letter: [%s]\n", to)
		}
		if to[1] < '1' || to[1] > '8' {
			fmt.Printf("bad target row digit: [%s]\n", to)
		}

		//b := game.history[len(game.history)-1]                        // will update a copy
		p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board

		if len(move) > 4 {
			// promotion
			promotion := move[4]
			kind := pieceKindFromLetter(rune(promotion))
			if kind != pieceNone {
				p = piece(b.turn<<3) + kind
			}
		}

		b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p) // put piece on board
		b.turn = colorInverse(b.turn)                           // switch color
		b.lastMove = fmt.Sprintf("%s %s", from, to)             // record move

		game.history = append(game.history, b) // append to history
	*/
}

func cmdReset(cmds []command, game *gameState, tokens []string) {
	game.loadFromString(builtinBoard)
}

func cmdSwitch(cmds []command, game *gameState, tokens []string) {
	b := &game.history[len(game.history)-1] // will update in place
	b.turn = colorInverse(b.turn)           // switch color
}

func cmdUndo(cmds []command, game *gameState, tokens []string) {
	if len(game.history) < 2 {
		return
	}
	game.history = game.history[:len(game.history)-1]
}
