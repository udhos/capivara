package main

import (
	"fmt"
	"strings"
)

func (g gameState) showFen() {
	fmt.Printf("fen: ")
	last := len(g.history) - 1
	b := g.history[last]

	// rows
	showFenRow(b, 7)
	for row := location(6); row >= 0; row-- {
		fmt.Print("/")
		showFenRow(b, row)
	}

	// turn
	if b.turn == 0 {
		fmt.Print(" w")
	} else {
		fmt.Print(" b")
	}

	// castling rights
	castling := ""
	if b.flags[0]&lostCastlingRight == 0 {
		castling += "K"
	}
	if b.flags[0]&lostCastlingRight == 0 {
		castling += "Q"
	}
	if b.flags[1]&lostCastlingRight == 0 {
		castling += "k"
	}
	if b.flags[1]&lostCastlingRight == 0 {
		castling += "q"
	}
	if castling == "" {
		fmt.Print(" -")
	} else {
		fmt.Print(" ", castling)
	}

	// FIXME - en passant target square
	fmt.Print(" -")

	// FIXME - Halfmove clock: This is the number of halfmoves since the last capture or pawn advance.
	fmt.Print(" 0")

	// Fullmove clock
	fmt.Print(" ", 1+(len(g.history)-1)/2)

	fmt.Println()
}

func showFenRow(b board, row location) {
	emptySquares := 0

	// first column
	loc := row*8 + 0
	p := b.square[loc]
	if p == pieceNone {
		emptySquares++
	} else {
		fmt.Print(fenLetter(p))
	}

	for col := location(1); col < 7; col++ {
		loc := row*8 + col
		p := b.square[loc]
		if p == pieceNone {
			emptySquares++
		} else {
			if emptySquares > 0 {
				fmt.Print(emptySquares)
				emptySquares = 0
			}
			fmt.Print(fenLetter(p))
		}
	}

	// last column
	loc = row*8 + 7
	p = b.square[loc]
	if p == pieceNone {
		emptySquares++
		fmt.Print(emptySquares)
	} else {
		if emptySquares > 0 {
			fmt.Print(emptySquares)
		}
		fmt.Print(fenLetter(p))
	}
}

func fenLetter(p piece) string {
	low := p.kindLetterLow()
	if p.color() == colorWhite {
		return strings.ToUpper(low)
	}
	return low
}
