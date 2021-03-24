package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type location uint8
type colorFlag uint32

const (
	lostCastlingLeft colorFlag = 1 << iota
	lostCastlingRight
)

type board struct {
	king          [2]location // king location
	square        [64]piece
	flags         [2]colorFlag
	turn          pieceColor
	materialValue [2]int
}

func (b *board) addPiece(i, j location, p piece) {

	b.delPiece(i, j)

	loc := i*8 + j
	//fmt.Printf("addPiece: %dx%d=%d color=%d kind=%s\n", i, j, loc, p.color(), p.kindLetter())
	b.square[loc] = p

	// record king position
	if p.kind() == whiteKing {
		b.king[p.color()] = loc
	}

	b.materialValue[p.color()] += p.materialValue() // piece material value enters board
}

func (b *board) delPiece(i, j location) piece {
	loc := i*8 + j
	p := b.square[loc]

	b.materialValue[p.color()] -= p.materialValue() // piece material value leaves board

	b.square[loc] = pieceNone

	return p
}

func (b board) getMaterialValue() float32 {
	wh := float32(b.materialValue[0])
	bl := float32(b.materialValue[1])
	return (wh + bl) / 100
}

type gameState struct {
	root board
}

func (g gameState) show() {
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Println("   -------------------------")
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d  |", row+1)
		for col := 0; col < 8; col++ {
			loc := row*8 + col
			piece := g.root.square[loc]
			piece.show()
			fmt.Print("|")
		}
		fmt.Printf("  %d\n", row+1)
		fmt.Println("   -------------------------")
	}
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Printf("turn: %s\n", g.root.turn.name())
	fmt.Printf("material: %v\n", g.root.getMaterialValue())
	fmt.Printf("white king: %dx%d material=%d\n", g.root.king[0]/8, g.root.king[0]%8, g.root.materialValue[0])
	fmt.Printf("black king: %dx%d material=%d\n", g.root.king[1]/8, g.root.king[1]%8, g.root.materialValue[1])
}

func (g *gameState) load(filename string) {
	fmt.Printf("loading: %s\n", filename)
	input, errOpen := os.Open(filename)
	if errOpen != nil {
		fmt.Printf("load: %s: %v\n", filename, errOpen)
		return
	}
	defer input.Close()

	reader := bufio.NewReader(input)

	var lineCount int

	for {
		lineCount++
		line, errRead := reader.ReadString('\n')
		switch errRead {
		case io.EOF:
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
				g.root.addPiece(location(row), location(col), piece(color<<3)+whitePawn)
			case c == 'R':
				g.root.addPiece(location(row), location(col), piece(color<<3)+whiteRook)
			case c == 'N':
				g.root.addPiece(location(row), location(col), piece(color<<3)+whiteKnight)
			case c == 'B':
				g.root.addPiece(location(row), location(col), piece(color<<3)+whiteBishop)
			case c == 'Q':
				g.root.addPiece(location(row), location(col), piece(color<<3)+whiteQueen)
			case c == 'K':
				g.root.addPiece(location(row), location(col), piece(color<<3)+whiteKing)
			}
		}
	}
}

func main() {
	game_loop()
}

func game_loop() {

	game := gameState{}

	game.load("board.txt")

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
				cmd.call(&game, tokens)
				continue LOOP
			}
		}

		fmt.Printf("bad command: %s\n", cmdPrefix)
	}
}

var tableCmd = []struct {
	name string
	call func(game *gameState, tokens []string)
}{
	{"load", cmdLoad},
	{"move", cmdMove},
}

func cmdLoad(game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: load filename\n")
		return
	}
	game.load(tokens[1])
}

func cmdMove(game *gameState, tokens []string) {
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

	p := game.root.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board

	game.root.addPiece(location(to[1]-'1'), location(to[0]-'a'), p) // put piece on board
}
