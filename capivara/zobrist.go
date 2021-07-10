package main

import (
	"fmt"
	"math/rand"
)

type zobristKey uint64

func (k zobristKey) String() string {
	return fmt.Sprintf("%016X", uint64(k))
}

// kind:
//  white 1 .. 6
//  black 9 .. 14
const pieceTypes = 14 // 0..13

var zobristBoard [64][pieceTypes]zobristKey
var zobristTurn [2]zobristKey
var zobristCastling [16]zobristKey
var zobristEnpassantCol [8]zobristKey
var zobristRand *rand.Rand

func zobristInit() {

	zobristRand = rand.New(rand.NewSource(20210709))

	for s := 0; s < 64; s++ {
		for p := 0; p < pieceTypes; p++ {
			zobristBoard[s][p] = zobristRandKey()
		}
	}
	zobristTurn[0] = zobristRandKey()
	zobristTurn[1] = zobristRandKey()
	for i := 0; i < 16; i++ {
		zobristCastling[i] = zobristRandKey()
	}
	for i := 0; i < 8; i++ {
		zobristEnpassantCol[i] = zobristRandKey()
	}
}

func zobristRandKey() zobristKey {
	return zobristKey(zobristRand.Uint64())
}

func (b *board) zobristInit() {
	b.zobristValue = zobristTurn[b.turn]

	b.zobristUpdateCastling()

	b.zobristUpdateEnPassant()

	for s := 0; s < 64; s++ {
		p := b.square[s]
		if p == pieceNone {
			continue
		}
		b.zobristUpdatePiece(s, p)
	}
}

func (b *board) zobristUpdateTurn() {
	b.zobristValue ^= zobristTurn[b.turn]
}

func (b *board) zobristUpdateCastling() {
	b.zobristValue ^= zobristCastling[0xF&(b.flags[0]|(b.flags[1]<<2))]
}

func (b *board) zobristUpdateEnPassant() {
	b.zobristValue ^= zobristEnpassantCol[7&b.lastMove.dst]
}

func (b *board) zobristUpdatePiece(loc int, p piece) {

	//p := b.square[loc]
	// kind:
	//  white 1 .. 6
	//  black 9 .. 14
	k0 := p.kind() - 1 // 0..13
	//fmt.Printf("zobristUpdatePiece: loc=%d kind=%d k0=%d\n", loc, p.kind(), k0)
	b.zobristValue ^= zobristBoard[loc][k0]
}
