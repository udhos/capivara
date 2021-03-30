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

func (p piece) materialValue() int {
	switch p {
	case pieceNone:
		return 0
	case whitePawn:
		return 100
	case blackPawn:
		return -100
	case whiteRook:
		return 500
	case whiteBishop:
		return 300
	case whiteKnight:
		return 300
	case blackRook:
		return -500
	case blackBishop:
		return -300
	case blackKnight:
		return -300
	case whiteQueen:
		return 900
	case blackQueen:
		return -900
	}
	return 0
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

func moveToStr(src, dst location) string {
	return fmt.Sprintf("%s%s", locToStr(src), locToStr(dst))
}

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
