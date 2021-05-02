package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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
	if b.flags[0]&lostCastlingLeft == 0 {
		castling += "Q"
	}
	if b.flags[1]&lostCastlingRight == 0 {
		castling += "k"
	}
	if b.flags[1]&lostCastlingLeft == 0 {
		castling += "q"
	}
	if castling == "" {
		fmt.Print(" -")
	} else {
		fmt.Print(" ", castling)
	}

	// En passant target square
	fmt.Print(" ", passantSquare(b))

	// FIXME - Halfmove clock: This is the number of halfmoves since the last capture or pawn advance.
	fmt.Print(" 0")

	// Fullmove clock
	fmt.Print(" ", 1+(len(g.history)-1)/2)

	fmt.Println()
}

func passantSquare(b board) string {
	if len(b.lastMove) == 4 {
		move := b.lastMove
		step := int64(move[3]) - int64(move[1])
		if abs(step) == 2 {
			// moved two squares
			dstCol := location(move[2] - 'a')
			dstRow := location(move[3] - '1')
			dstLoc := dstRow*8 + dstCol
			p := b.square[dstLoc]
			if p.kind() == whitePawn {
				// it is pawn
				targetRow := dstRow - location(colorToSignal(p.color()))
				return coordToStr(targetRow, dstCol)
			}
		}
	}
	return "-"
}

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
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

func fenParse(fen []string) (board, error) {
	b := board{}

	// drop castling rights
	b.flags[colorWhite] |= lostCastlingLeft | lostCastlingRight // disable castling for white
	b.flags[colorBlack] |= lostCastlingLeft | lostCastlingRight // disable castling for black

	fields := len(fen)

	// parse pieces

	if fields < 1 {
		return b, errors.New("missing FEN pieces") // no pieces
	}

	rows := strings.FieldsFunc(fen[0], func(r rune) bool { return r == '/' })
	for r, codeRow := range rows {
		row := 7 - r
		col := 0
		for _, codeCol := range codeRow {
			if codeCol >= '1' && codeCol <= '8' {
				col += int(codeCol) - '0'
				continue
			}
			kind := pieceKindFromLetter(codeCol)
			color := colorBlack
			if unicode.IsUpper(codeCol) {
				color = colorWhite
			}
			b.addPiece(location(row), location(col), piece(color<<3)+kind)
			col++
		}
	}

	// parse turn

	if fields < 2 {
		return b, nil // no turn
	}

	if fen[1] != "w" {
		b.turn = colorBlack
	}

	// parse castling rights

	if fields < 3 {
		return b, nil // no castling rights
	}

	for _, l := range fen[2] {
		switch l {
		case 'K':
			b.flags[colorWhite] &= ^lostCastlingRight
		case 'Q':
			b.flags[colorWhite] &= ^lostCastlingLeft
		case 'k':
			b.flags[colorBlack] &= ^lostCastlingRight
		case 'q':
			b.flags[colorBlack] &= ^lostCastlingLeft
		}
	}

	return b, nil
}
