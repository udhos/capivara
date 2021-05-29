package main

import (
	"fmt"
	"unicode"
)

type pieceColor uint8

func (p pieceColor) name() string {
	switch p {
	case colorWhite:
		return "white"
	}
	return "black"
}

const (
	colorWhite pieceColor = iota
	colorBlack
)

type piece uint8

const (
	pieceNone piece = iota
	whiteKing
	whiteQueen
	whiteRook
	whiteBishop
	whiteKnight
	whitePawn
)

const (
	blackKing piece = iota + whiteKing + 8
	blackQueen
	blackRook
	blackBishop
	blackKnight
	blackPawn
)

// pieceColor: 0=white 1=black
func (p piece) color() pieceColor {
	return pieceColor(p >> 3)
}

func (p piece) kind() piece {
	return piece(p & 7)
}

func (p piece) kindLetter() string {
	switch p.kind() {
	case whiteKing:
		return "K"
	case whiteQueen:
		return "Q"
	case whiteRook:
		return "R"
	case whiteBishop:
		return "B"
	case whiteKnight:
		return "N"
	case whitePawn:
		return "p"
	}
	return "?"
}

func (p piece) kindLetterLow() string {
	switch p.kind() {
	case whitePawn:
		return "p"
	case whiteRook:
		return "r"
	case whiteBishop:
		return "b"
	case whiteKnight:
		return "n"
	case whiteKing:
		return "k"
	case whiteQueen:
		return "q"
	}
	return "?"
}

func pieceKindFromLetter(letter rune) piece {
	switch unicode.ToLower(letter) {
	case 'k':
		return whiteKing
	case 'q':
		return whiteQueen
	case 'r':
		return whiteRook
	case 'b':
		return whiteBishop
	case 'n':
		return whiteKnight
	case 'p':
		return whitePawn
	}
	return pieceNone
}

func (p piece) materialValue(loc location) int16 {
	switch p {
	case pieceNone:
		return 0
	case whitePawn:
		return p.piecePlusPosition(100, loc)
	case blackPawn:
		return -p.piecePlusPosition(100, loc)
	case whiteRook:
		return p.piecePlusPosition(500, loc)
	case whiteBishop:
		return p.piecePlusPosition(300, loc)
	case whiteKnight:
		return p.piecePlusPosition(250, loc)
	case blackRook:
		return -p.piecePlusPosition(500, loc)
	case blackBishop:
		return -p.piecePlusPosition(300, loc)
	case blackKnight:
		return -p.piecePlusPosition(250, loc)
	case whiteQueen:
		return p.piecePlusPosition(900, loc)
	case blackQueen:
		return -p.piecePlusPosition(900, loc)
	}
	return 0
}

func (p piece) piecePlusPosition(value int16, loc location) int16 {
	return value + positionTable[p.kind()-1][loc]
}

func (p piece) show() {
	if p == pieceNone {
		fmt.Print("  ")
		return
	}
	color := p.color()
	if color == colorWhite {
		fmt.Print("*")
	} else {
		fmt.Print(".")
	}
	fmt.Print(p.kindLetter())
}

func coordToStr(row, col location) string {
	return fmt.Sprintf("%c%d", col+'a', row+1)
}

func locToStr(loc location) string {
	return coordToStr(loc/8, loc%8)
}

/*
func moveToStr(src, dst location, p piece) string {
	if p == pieceNone {
		return fmt.Sprintf("%s%s", locToStr(src), locToStr(dst))
	}
	return fmt.Sprintf("%s%s%s", locToStr(src), locToStr(dst), p.kindLetterLow())
}
*/

// white=0 -> signal=1
// black=1 -> signal=-1
func colorToSignal(color pieceColor) int {
	return 1 - int(color)*2
}

// white=0 -> black=1
// black=1 -> white=0
func colorInverse(color pieceColor) pieceColor {
	return 1 - color
}
