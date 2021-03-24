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
	king   [2]int // king location
	square [64]piece
	flags  [2]colorFlag
}

func (b *board) addPiece(i, j location, p piece) {
	loc := i*8 + j
	//fmt.Printf("addPiece: %dx%d=%d color=%d kind=%s\n", i, j, loc, p.color(), p.kindLetter())
	b.square[loc] = p
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
		cmd := tokens[0]
		if strings.HasPrefix("load", cmd) {
			if len(tokens) < 2 {
				fmt.Printf("usage: load filename\n")
				continue
			}
			game.load(tokens[1])
			continue
		}

		fmt.Printf("bad command: %s\n", cmd)
	}
}
