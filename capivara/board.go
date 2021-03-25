package main

import "fmt"

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
	lastMove      string
}

func (b *board) addPiece(i, j location, p piece) {
	loc := i*8 + j
	b.addPieceLoc(loc, p)
}

func (b *board) delPiece(i, j location) piece {
	loc := i*8 + j
	return b.delPieceLoc(loc)
}

func (b *board) addPieceLoc(loc location, p piece) {
	b.delPieceLoc(loc)

	b.square[loc] = p

	// record king position
	if p.kind() == whiteKing {
		b.king[p.color()] = loc
	}

	b.materialValue[p.color()] += p.materialValue() // piece material value enters board
}

func (b *board) delPieceLoc(loc location) piece {
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

func (b board) generateChildren(children []board) []board {
	//var n int
	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			continue
		}
		children = b.generateChildrenPiece(children, loc, p)
	}
	return children
}

func (b board) generateChildrenPiece(children []board, loc location, p piece) []board {
	//var n int
	i, j := loc/8, loc%8
	kind := p.kind()
	color := p.color()
	signal := colorToSignal(color)              // 0=>1 1=>-1
	lastRow := location(7 - 7*int(color))       // 0=>7 1=>0
	firstRow := location(7*int(color) + signal) // 0=>1 1=>6
	switch kind {
	case whitePawn: // white + black
		// can move one up/down?
		if i != lastRow {
			dstLoc := (i+location(signal))*8 + j
			dstP := b.square[dstLoc]
			if dstP == pieceNone {
				// position is free
				child := b
				pp := child.delPieceLoc(loc)                                           // take piece from board
				child.addPieceLoc(dstLoc, pp)                                          // put piece on board
				child.turn = colorInverse(b.turn)                                      // switch color
				child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstLoc)) // record move
				children = append(children, child)                                     // append to children
			}
		}

		// can move two up/down?
		if i == firstRow {
			secondRow := firstRow + location(signal)
			dstRow := secondRow + location(signal)
			secondRowLoc := secondRow*8 + j
			dstRowLoc := dstRow*8 + j
			secondP := b.square[secondRowLoc]
			dstP := b.square[dstRowLoc]
			if secondP == pieceNone && dstP == pieceNone {
				// free to move
				child := b
				pp := child.delPieceLoc(loc)                                              // take piece from board
				child.addPieceLoc(dstRowLoc, pp)                                          // put piece on board
				child.turn = colorInverse(b.turn)                                         // switch color
				child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstRowLoc)) // record move
				children = append(children, child)                                        // append to children
			}
		}
	}

	return children
}