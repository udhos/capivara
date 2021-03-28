package main

type location int8
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
	i, j := int(loc)/8, int(loc)%8
	kind := p.kind()
	color := p.color()
	signal := colorToSignal(color)    // 0=>1 1=>-1
	lastRow := 7 - 7*int(color)       // 0=>7 1=>0
	firstRow := 7*int(color) + signal // 0=>1 1=>6
	switch kind {
	case whitePawn: // white + black
		// can move one up/down?
		{
			dstRow := i + signal
			dstLoc := dstRow*8 + j
			dstP := b.square[dstLoc]
			if dstP == pieceNone {
				// position is free
				if dstRow == lastRow {
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					children = b.recordMoveIfValid(children, loc, location(dstLoc))
				}
			}
		}

		// can move two up/down?
		if i == firstRow {
			secondRow := firstRow + signal
			dstRow := secondRow + signal
			secondRowLoc := secondRow*8 + j
			dstRowLoc := dstRow*8 + j
			secondP := b.square[secondRowLoc]
			dstP := b.square[dstRowLoc]
			if secondP == pieceNone && dstP == pieceNone {
				// free to move
				children = b.recordMoveIfValid(children, loc, location(dstRowLoc))
			}
		}

		// capture left?
		if j > 0 && i > 0 && i < 7 {
			dstRow := i + signal
			dstLoc := dstRow*8 + j - 1
			dstP := b.square[dstLoc]
			if dstP != pieceNone && dstP.color() != color {
				// free to capture
				if dstRow == lastRow {
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					children = b.recordMoveIfValid(children, loc, location(dstLoc))
				}
			}
		}

		// capture right?
		if j < 7 && i > 0 && i < 7 {
			dstRow := i + signal
			dstLoc := dstRow*8 + j + 1
			dstP := b.square[dstLoc]
			if dstP != pieceNone && dstP.color() != color {
				// free to capture
				if dstRow == lastRow {
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					children = b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					children = b.recordMoveIfValid(children, loc, location(dstLoc))
				}
			}
		}

	case whiteQueen: // white + black
		children = b.generateSliding(children, loc, 0, 1)
		children = b.generateSliding(children, loc, 1, 1)
		children = b.generateSliding(children, loc, 1, 0)
		children = b.generateSliding(children, loc, 1, -1)
		children = b.generateSliding(children, loc, 0, -1)
		children = b.generateSliding(children, loc, -1, -1)
		children = b.generateSliding(children, loc, -1, 0)
		children = b.generateSliding(children, loc, -1, 1)
	}

	return children
}

func (b board) generateSliding(children []board, src, incRow, incCol location) []board {
	dstRow := src / 8
	dstCol := src % 8
	for {
		dstRow += incRow
		dstCol += incCol
		if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
			break // out of board
		}
		dstLoc := dstRow*8 + dstCol
		dstP := b.square[dstLoc]
		if dstP == pieceNone {
			// empty square
			children = b.recordMoveIfValid(children, src, dstLoc)
			continue
		}
		if dstP.color() == b.turn {
			break // blocked by same color piece
		}
		// capture opponent piece
		children = b.recordMoveIfValid(children, src, dstLoc)
		break
	}

	return children
}

func (b board) recordMoveIfValid(children []board, src, dst location) []board {
	child := b                           // copy board
	p := child.delPieceLoc(src)          // take piece from board
	child.addPieceLoc(dst, p)            // put piece on board
	child.turn = colorInverse(b.turn)    // switch color
	child.lastMove = moveToStr(src, dst) // record move

	if child.otherKingInCheck() {
		return children // drop invalid move
	}

	children = append(children, child) // append valid move to children
	return children
}

func (b board) recordPromotionIfValid(children []board, src, dst location, p piece) []board {
	child := b                           // copy board
	child.delPieceLoc(src)               // take pawn from board
	child.addPieceLoc(dst, p)            // put new piece on board
	child.turn = colorInverse(b.turn)    // switch color
	child.lastMove = moveToStr(src, dst) // record move

	if child.otherKingInCheck() {
		return children // drop invalid move
	}

	children = append(children, child) // append valid move to children
	return children
}

func (b board) otherKingInCheck() bool {
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

func (b board) kingInCheck() bool {
	kingLoc := b.king[b.turn]
	kingPiece := b.square[kingLoc]

	// any piece attacks king?
	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() == b.turn {
			continue
		}
		if b.pieceAttacks(p, loc, kingPiece, kingLoc) {
			return true // king is in check
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
	case whiteQueen: // white + black
		if b.slidingAttack(srcLoc, dstLoc, 0, 1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 1, 1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 1, 0) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 1, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 0, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, 0) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, 1) {
			return true
		}
	}

	return false
}

func (b board) slidingAttack(src, target, incRow, incCol location) bool {
	dstRow := src / 8
	dstCol := src % 8
	for {
		dstRow += incRow
		dstCol += incCol
		if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
			break // out of board
		}
		dstLoc := dstRow*8 + dstCol
		if dstLoc == target {
			return true // found
		}
		dstP := b.square[dstLoc]
		if dstP != pieceNone {
			break // blocked by some piece
		}
	}
	return false
}
