package main

import "fmt"

type piece uint8
type location uint8

const (
	whiteKing = piece(iota)
	whiteQueen
	whiteRook
	whiteBishop
	whiteKnight
	whitePawn
)

const (
	blackKing = piece(iota + whiteKing + 8)
	blackQueen
	blackRook
	blackBishop
	blackKnight
	blackPawn
)

// color: 0=white 1=black
func (p piece) color() uint8 {
	return uint8(p >> 3)
}

type board struct {
	piece  [2][]location // list o locations with pieces
	square [64]piece
}

func (b *board) addPiece(i, j location, p piece) {
	loc := i*8 + j
	b.square[loc] = p
	color := p.color()
	b.piece[color] = append(b.piece[color], loc)
}

func main() {
	fmt.Println("whitePawn:", whitePawn, " color:", whitePawn.color())
	fmt.Println("blackPawn:", blackPawn, " color:", blackPawn.color())
}
