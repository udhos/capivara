package main

import (
	"fmt"
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

func (p piece) materialValue() int {
	switch p {
	case whiteQueen:
		return 900
	case whiteRook:
		return 500
	case whiteBishop:
		return 300
	case whiteKnight:
		return 300
	case whitePawn:
		return 100
	case blackQueen:
		return -900
	case blackRook:
		return -500
	case blackBishop:
		return -300
	case blackKnight:
		return -300
	case blackPawn:
		return -100
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
