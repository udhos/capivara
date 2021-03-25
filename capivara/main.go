package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
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
	fmt.Printf("white king: %s material=%d\n", locToStr(b.king[0]), b.materialValue[0])
	fmt.Printf("black king: %s material=%d\n", locToStr(b.king[1]), b.materialValue[1])
	fmt.Printf("history %d moves: ", len(g.history))
	for _, m := range g.history {
		fmt.Printf("(%s)", m.lastMove)
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

func main() {
	gameLoop()
}

func gameLoop() {

	game := newGame()
	game.loadFromString(builtinBoard)

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
	fmt.Print("available commands:")
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
	if len(tokens) < 3 {
		fmt.Printf("usage: move from to\n")
		return
	}
	from := strings.ToLower(tokens[1])
	to := strings.ToLower(tokens[2])

	// from
	if len(from) != 2 {
		fmt.Printf("bad source format: [%s]", from)
		return
	}
	if from[0] < 'a' || from[0] > 'h' {
		fmt.Printf("bad source column letter: [%s]", from)
	}
	if from[1] < '1' || from[1] > '8' {
		fmt.Printf("bad source row digit: [%s]", from)
	}

	// to
	if len(to) != 2 {
		fmt.Printf("bad target format: [%s]", to)
		return
	}
	if to[0] < 'a' || to[0] > 'h' {
		fmt.Printf("bad target column letter: [%s]", to)
	}
	if to[1] < '1' || to[1] > '8' {
		fmt.Printf("bad target row digit: [%s]", to)
	}

	b := &game.history[len(game.history)-1]                       // will update in-place
	p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board
	b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p)       // put piece on board
}

func cmdNegamax(cmds []command, game *gameState, tokens []string) {
	last := len(game.history) - 1
	b := game.history[last]
	rootNegamax(b, 2)
}

func rootNegamax(b board, depth int) float32 {
	terminalNode := false // checkmate or draw?
	if depth == 0 || terminalNode {
		return b.getMaterialValue()
	}
	var max float32 = -1000.0

	children := []board{}
	children = b.generateChildren(children)
	if len(children) == 0 {
		fmt.Printf("rootNegamax: %s %s UGH NO CHILDREN CHECK-MATE?\n", b.turn.name(), b.lastMove)
	}
	for i, child := range children {
		score := -negamax(child, depth-1)
		fmt.Printf("rootNegamax: child: %d: %s %s = score=%v max=%v\n", i, child.turn.name(), child.lastMove, score, max)
		if score > max {
			max = score
			fmt.Printf("rootNegamax: child: %d: %s %s = score=%v max=%v BEST\n", i, child.turn.name(), child.lastMove, score, max)
		}
	}
	return max
}

func negamax(b board, depth int) float32 {
	terminalNode := false // checkmate or draw?
	if depth == 0 || terminalNode {
		fmt.Printf("negamax: %s %s depth=%d value=%v MATERIAL\n", b.turn.name(), b.lastMove, depth, b.getMaterialValue())
		return b.getMaterialValue()
	}
	var value float32 = -1000.0

	children := []board{}
	children = b.generateChildren(children)
	if len(children) == 0 {
		fmt.Printf("negamax: child: %s %s depth=%d value=%v UGH NO CHILDREN CHECK-MATE?\n", b.turn.name(), b.lastMove, depth, value)
	}
	for _, child := range children {
		value = max(value, -negamax(child, depth-1))
	}
	fmt.Printf("negamax: child: %s %s depth=%d value=%v\n", b.turn.name(), b.lastMove, depth, value)
	return value
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func cmdPlay(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 3 {
		fmt.Printf("usage: play from to\n")
		return
	}
	from := strings.ToLower(tokens[1])
	to := strings.ToLower(tokens[2])

	// from
	if len(from) != 2 {
		fmt.Printf("bad source format: [%s]", from)
		return
	}
	if from[0] < 'a' || from[0] > 'h' {
		fmt.Printf("bad source column letter: [%s]", from)
	}
	if from[1] < '1' || from[1] > '8' {
		fmt.Printf("bad source row digit: [%s]", from)
	}

	// to
	if len(to) != 2 {
		fmt.Printf("bad target format: [%s]", to)
		return
	}
	if to[0] < 'a' || to[0] > 'h' {
		fmt.Printf("bad target column letter: [%s]", to)
	}
	if to[1] < '1' || to[1] > '8' {
		fmt.Printf("bad target row digit: [%s]", to)
	}

	b := game.history[len(game.history)-1]                        // will update a copy
	p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board
	b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p)       // put piece on board
	b.turn = colorInverse(b.turn)                                 // switch color
	b.lastMove = fmt.Sprintf("%s %s", from, to)                   // record move

	game.history = append(game.history, b) // append to history
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
