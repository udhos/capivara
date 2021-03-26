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
				/*
					child := b
					pp := child.delPieceLoc(loc)                                           // take piece from board
					child.addPieceLoc(dstLoc, pp)                                          // put piece on board
					child.turn = colorInverse(b.turn)                                      // switch color
					child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstLoc)) // record move
					children = append(children, child)                                     // append to children
				*/
				children = b.recordMoveIfValid(children, loc, dstLoc)
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
				/*
					child := b
					pp := child.delPieceLoc(loc)                                              // take piece from board
					child.addPieceLoc(dstRowLoc, pp)                                          // put piece on board
					child.turn = colorInverse(b.turn)                                         // switch color
					child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstRowLoc)) // record move
					children = append(children, child)                                        // append to children
				*/
				children = b.recordMoveIfValid(children, loc, dstRowLoc)
			}
		}

		// capture left?
		if j > 0 {
			dstRow := i + location(signal)
			dstLoc := dstRow*8 + j - 1
			dstP := b.square[dstLoc]
			if dstP != pieceNone && dstP.color() != color {
				// free to capture
				/*
					child := b
					pp := child.delPieceLoc(loc)                                           // take piece from board
					child.addPieceLoc(dstLoc, pp)                                          // put piece on board
					child.turn = colorInverse(b.turn)                                      // switch color
					child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstLoc)) // record move
					children = append(children, child)                                     // append to children
				*/
				children = b.recordMoveIfValid(children, loc, dstLoc)
			}
		}

		// capture right?
		if j < 7 {
			dstRow := i + location(signal)
			dstLoc := dstRow*8 + j + 1
			dstP := b.square[dstLoc]
			if dstP != pieceNone && dstP.color() != color {
				// free to capture
				/*
					child := b
					pp := child.delPieceLoc(loc)                                           // take piece from board
					child.addPieceLoc(dstLoc, pp)                                          // put piece on board
					child.turn = colorInverse(b.turn)                                      // switch color
					child.lastMove = fmt.Sprintf("%s %s", locToStr(loc), locToStr(dstLoc)) // record move
					children = append(children, child)                                     // append to children
				*/
				children = b.recordMoveIfValid(children, loc, dstLoc)
			}
		}
	}

	return children
}

func (b board) recordMoveIfValid(children []board, src, dst location) []board {
	child := b                                                          // copy board
	p := child.delPieceLoc(src)                                         // take piece from board
	child.addPieceLoc(dst, p)                                           // put piece on board
	child.turn = colorInverse(b.turn)                                   // switch color
	child.lastMove = fmt.Sprintf("%s %s", locToStr(src), locToStr(dst)) // record move

	if child.invalid() {
		return children
	}

	children = append(children, child) // append to children
	return children
}

func (b board) invalid() bool {
	otherKingColor := colorInverse(b.turn)
	otherKingLoc := b.king[otherKingColor]
	otherKingPiece := b.square[otherKingLoc]

	// any piece attacks other king?
	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			continue
		}
		if b.pieceAttacks(p, loc, otherKingPiece, otherKingLoc) {
			return true // other king is in check
		}
	}

	return false
}

func (b board) pieceAttacks(srcPiece piece, srcLoc location, dstPiece piece, dstLoc location) bool {

	srcRow, srcCol := srcLoc/8, srcLoc%8
	dstRow, dstCol := dstLoc/8, dstLoc%8

	srcKind := srcPiece.kind()
	srcColor := srcPiece.color()
	srcSignal := colorToSignal(srcColor) // 0=>1 1=>-1

	switch srcKind {
	case whitePawn: // white + black
		if int(srcRow)+srcSignal == int(dstRow) {
			return (dstCol == srcCol-1) || (dstCol == srcCol+1)
		}
	}

	return false
}
