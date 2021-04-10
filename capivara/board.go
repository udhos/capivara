package main

type location int8
type colorFlag uint8

const (
	lostCastlingLeft colorFlag = 1 << iota
	lostCastlingRight
)

type board struct {
	king          [2]location // king location
	square        [64]piece
	flags         [2]colorFlag
	turn          pieceColor
	materialValue [2]int16
	lastMove      string
}

func (b *board) validMove(move string) bool {
	for _, c := range b.generateChildren([]board{}) {
		if move == c.lastMove {
			return true
		}
	}
	return false
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
	b.materialValue[p.color()] += p.materialValue() // piece material value enters board
	if p.kind() == whiteKing {
		// record king new position
		b.king[p.color()] = loc
	}
}

func (b *board) delPieceLoc(loc location) piece {
	p := b.square[loc]
	b.materialValue[p.color()] -= p.materialValue() // piece material value leaves board
	b.square[loc] = pieceNone

	/*
		switch p.kind() {
		case whiteRook:
			// rook moved, then disable castling
			firstRow := 7 * location(p.color()) // 0=>0 1=>7
			row := loc / 8
			if row == firstRow {
				// rook in initial row
				col := loc % 8
				switch col {
				case 0: // left rook
					b.flags[b.turn] |= lostCastlingLeft
				case 7: // right rook
					b.flags[b.turn] |= lostCastlingRight
				}
			}
		case whiteKing:
			// king moved, then disable castling
			b.flags[b.turn] |= lostCastlingLeft | lostCastlingRight
		}
	*/

	return p
}

func (b board) getMaterialValue() float32 {
	wh := float32(b.materialValue[0])
	bl := float32(b.materialValue[1])
	return (wh + bl) / 100
}

func (b board) generateChildren(children []board) []board {

	// generate en passant captures

	lastMove := b.lastMove
	if lastMove != "" {
		trgColor := b.turn
		trgSignal := colorToSignal(trgColor)

		trgSrcRow := int(lastMove[1]) - '0'
		trgDstRow := int(lastMove[3]) - '0'

		trgRowFrom := 7*int(trgColor) + trgSignal // 0=>1 1=>6
		trgRowTo := trgRowFrom + 2*trgSignal      // 1=>3 7=>5

		// from 2nd to 4th ?
		if trgSrcRow == trgRowFrom && trgDstRow == trgRowTo {
			trgSrcCol := int(lastMove[0]) - 'a'
			trgDstCol := int(lastMove[2]) - 'a'

			// same column?
			if trgSrcCol == trgDstCol {
				trgDstLoc := trgDstCol + 8*trgDstRow
				trgKind := b.square[trgDstLoc].kind()
				if trgKind == whitePawn {
					// it is pawn
					if trgDstCol > 0 {
						// might be captured from left
					}
					if trgDstCol < 7 {
						// might be captured from right
					}
				}
			}
		}
	}

	// scan pieces
	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone || p.color() != b.turn {
			continue
		}
		children = b.generateChildrenPiece(children, loc, p)
	}

	// generate castling

	if b.flags[b.turn]&lostCastlingLeft == 0 {
		// castling left
		firstRow8 := 8 * 7 * location(b.turn) // 0=>0 1=>7
		colB := firstRow8 + 1
		colC := firstRow8 + 2
		colD := firstRow8 + 3
		if b.square[colB] == pieceNone && b.square[colC] == pieceNone && b.square[colD] == pieceNone {
			// squares are free
			colE := firstRow8 + 4 // king
			if !b.anyPieceAttacks(colB) && !b.anyPieceAttacks(colC) && !b.anyPieceAttacks(colD) && !b.anyPieceAttacks(colE) {
				children = b.generateCastlingLeft(children)
			}
		}
	}
	if b.flags[b.turn]&lostCastlingRight == 0 {
		// castling right
		firstRow8 := 8 * 7 * location(b.turn) // 0=>0 1=>7
		colF := firstRow8 + 5
		colG := firstRow8 + 6
		if b.square[colF] == pieceNone && b.square[colG] == pieceNone {
			// squares are free
			colE := firstRow8 + 4 // king
			if !b.anyPieceAttacks(colE) && !b.anyPieceAttacks(colF) && !b.anyPieceAttacks(colG) {
				children = b.generateCastlingRight(children)
			}
		}
	}

	return children
}

func (b board) generateCastlingLeft(children []board) []board {

	row := 7 * location(b.turn)
	row8 := 8 * row
	kingSrc := row8 + 4 // E
	kingDst := row8 + 2 // C
	rookSrc := row8     // A
	rookDst := row8 + 3 // D

	child := b // copy board

	// move
	child.square[kingDst] = child.square[kingSrc]
	child.square[rookDst] = child.square[rookSrc]
	child.square[kingSrc] = pieceNone
	child.square[rookSrc] = pieceNone

	// disable castling
	child.flags[child.turn] |= lostCastlingLeft | lostCastlingRight

	child.turn = colorInverse(b.turn)                       // switch color
	child.lastMove = moveToStr(kingSrc, kingDst, pieceNone) // record move

	return b.recordIfValid(children, child)
}

func (b board) generateCastlingRight(children []board) []board {
	row := 7 * location(b.turn)
	row8 := 8 * row
	kingSrc := row8 + 4 // E
	kingDst := row8 + 6 // G
	rookSrc := row8 + 7 // H
	rookDst := row8 + 5 // F

	child := b // copy board

	// move
	child.square[kingDst] = child.square[kingSrc]
	child.square[rookDst] = child.square[rookSrc]
	child.square[kingSrc] = pieceNone
	child.square[rookSrc] = pieceNone

	// disable castling
	child.flags[child.turn] |= lostCastlingLeft | lostCastlingRight

	child.turn = colorInverse(b.turn)                       // switch color
	child.lastMove = moveToStr(kingSrc, kingDst, pieceNone) // record move

	return b.recordIfValid(children, child)
}

func (b board) generateChildrenPiece(children []board, loc location, p piece) []board {
	kind := p.kind()
	switch kind {
	case whitePawn: // white + black
		i, j := int(loc)/8, int(loc)%8
		color := p.color()
		signal := colorToSignal(color)    // 0=>1 1=>-1
		lastRow := 7 - 7*int(color)       // 0=>7 1=>0
		firstRow := 7*int(color) + signal // 0=>1 1=>6

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

	case whiteRook: // white + black
		children = b.generateSlidingRook(children, loc, 0, 1)
		children = b.generateSlidingRook(children, loc, 1, 0)
		children = b.generateSlidingRook(children, loc, 0, -1)
		children = b.generateSlidingRook(children, loc, -1, 0)

	case whiteBishop: // white + black
		children = b.generateSliding(children, loc, 1, 1)
		children = b.generateSliding(children, loc, 1, -1)
		children = b.generateSliding(children, loc, -1, -1)
		children = b.generateSliding(children, loc, -1, 1)

	case whiteKing: // white + black
		children = b.generateRelativeKing(children, loc, 0, 1)
		children = b.generateRelativeKing(children, loc, 1, 1)
		children = b.generateRelativeKing(children, loc, 1, 0)
		children = b.generateRelativeKing(children, loc, 1, -1)
		children = b.generateRelativeKing(children, loc, 0, -1)
		children = b.generateRelativeKing(children, loc, -1, -1)
		children = b.generateRelativeKing(children, loc, -1, 0)
		children = b.generateRelativeKing(children, loc, -1, 1)

	case whiteKnight: // white + black
		children = b.generateRelative(children, loc, -1, 2)
		children = b.generateRelative(children, loc, 1, 2)
		children = b.generateRelative(children, loc, 2, -1)
		children = b.generateRelative(children, loc, 2, 1)
		children = b.generateRelative(children, loc, -1, -2)
		children = b.generateRelative(children, loc, 1, -2)
		children = b.generateRelative(children, loc, -2, -1)
		children = b.generateRelative(children, loc, -2, 1)
	}

	return children
}

func (b board) recordIfValid(children []board, child board) []board {
	if child.otherKingInCheck() {
		return children // drop invalid move 'child'
	}
	return append(children, child) // record
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
		if dstP.color() != b.turn {
			// capture opponent piece
			children = b.recordMoveIfValid(children, src, dstLoc)
		}
		break
	}

	return children
}

func (b board) generateSlidingRook(children []board, src, incRow, incCol location) []board {
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
			children = b.recordMoveIfValidRook(children, src, dstLoc)
			continue
		}
		if dstP.color() != b.turn {
			// capture opponent piece
			children = b.recordMoveIfValidRook(children, src, dstLoc)
		}
		break
	}

	return children
}

func (b board) generateRelative(children []board, src, incRow, incCol location) []board {
	dstRow := src / 8
	dstCol := src % 8

	dstRow += incRow
	dstCol += incCol
	if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
		return children // out of board
	}
	dstLoc := dstRow*8 + dstCol
	dstP := b.square[dstLoc]
	if dstP == pieceNone {
		// empty square
		children = b.recordMoveIfValid(children, src, dstLoc)
		return children
	}
	if dstP.color() != b.turn {
		// capture opponent piece
		children = b.recordMoveIfValid(children, src, dstLoc)
		return children
	}

	// blocked by same color piece

	return children
}

func (b board) generateRelativeKing(children []board, src, incRow, incCol location) []board {
	dstRow := src / 8
	dstCol := src % 8

	dstRow += incRow
	dstCol += incCol
	if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
		return children // out of board
	}
	dstLoc := dstRow*8 + dstCol
	dstP := b.square[dstLoc]
	if dstP == pieceNone {
		// empty square
		children = b.recordMoveIfValidKing(children, src, dstLoc)
		return children
	}
	if dstP.color() != b.turn {
		// capture opponent piece
		children = b.recordMoveIfValidKing(children, src, dstLoc)
		return children
	}

	// blocked by same color piece

	return children
}

func (b board) newChild(src, dst location) (board, piece) {
	child := b                                      // copy board
	p := child.delPieceLoc(src)                     // take piece from board
	child.addPieceLoc(dst, p)                       // put piece on board
	child.turn = colorInverse(b.turn)               // switch color
	child.lastMove = moveToStr(src, dst, pieceNone) // record move
	return child, p
}

func (b board) recordMoveIfValid(children []board, src, dst location) []board {
	child, _ := b.newChild(src, dst)
	return b.recordIfValid(children, child)
}

func (b board) recordMoveIfValidKing(children []board, src, dst location) []board {
	child, p := b.newChild(src, dst)

	// king moved, then disable castling
	child.flags[p.color()] |= lostCastlingLeft | lostCastlingRight

	return b.recordIfValid(children, child)
}

func (b board) recordMoveIfValidRook(children []board, src, dst location) []board {
	child, p := b.newChild(src, dst)

	// rook moved, then disable castling
	firstRow := 7 * location(p.color()) // 0=>0 1=>7
	row := src / 8
	if row == firstRow {
		// rook in initial row
		col := src % 8
		switch col {
		case 0: // left rook
			child.flags[p.color()] |= lostCastlingLeft
		case 7: // right rook
			child.flags[p.color()] |= lostCastlingRight
		}
	}

	return b.recordIfValid(children, child)
}

func (b board) recordPromotionIfValid(children []board, src, dst location, p piece) []board {
	child := b                              // copy board
	child.delPieceLoc(src)                  // take pawn from board
	child.addPieceLoc(dst, p)               // put new piece on board
	child.turn = colorInverse(b.turn)       // switch color
	child.lastMove = moveToStr(src, dst, p) // record move

	return b.recordIfValid(children, child)
}

func (b board) otherKingInCheck() bool {
	otherKingColor := colorInverse(b.turn)
	otherKingLoc := b.king[otherKingColor]

	// any piece attacks other king?
	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone || p.color() != b.turn {
			continue
		}
		if b.pieceAttacks(p, loc, otherKingLoc) {
			return true // other king is in check
		}
	}

	return false
}

func (b board) kingInCheck() bool {
	return b.anyPieceAttacks(b.king[b.turn])
}

// any piece attacks square?
func (b board) anyPieceAttacks(target location) bool {

	for loc := location(0); loc < location(64); loc++ {
		p := b.square[loc]
		if p == pieceNone || p.color() == b.turn {
			continue
		}
		if b.pieceAttacks(p, loc, target) {
			return true // square attacked
		}
	}

	return false // square not attacked
}

func (b board) pieceAttacks(srcPiece piece, srcLoc, dstLoc location) bool {

	srcRow, srcCol := srcLoc/8, srcLoc%8
	dstRow, dstCol := dstLoc/8, dstLoc%8
	srcKind := srcPiece.kind()

	switch srcKind {
	case whitePawn: // white + black
		srcColor := srcPiece.color()
		srcSignal := colorToSignal(srcColor) // 0=>1 1=>-1
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

	case whiteRook: // white + black
		if b.slidingAttack(srcLoc, dstLoc, 0, 1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 1, 0) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 0, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, 0) {
			return true
		}

	case whiteBishop: // white + black
		if b.slidingAttack(srcLoc, dstLoc, 1, 1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, 1, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, -1) {
			return true
		}
		if b.slidingAttack(srcLoc, dstLoc, -1, 1) {
			return true
		}

	case whiteKing: // white + black
		if b.relativeAttack(srcLoc, dstLoc, 0, 1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 1, 1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 1, 0) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 1, -1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 0, -1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -1, -1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -1, 0) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -1, 1) {
			return true
		}

	case whiteKnight: // white + black
		if b.relativeAttack(srcLoc, dstLoc, -1, 2) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 1, 2) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 2, -1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 2, 1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -1, -2) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, 1, -2) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -2, -1) {
			return true
		}
		if b.relativeAttack(srcLoc, dstLoc, -2, 1) {
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

func (b board) relativeAttack(src, target, incRow, incCol location) bool {
	dstRow := src / 8
	dstCol := src % 8

	dstRow += incRow
	dstCol += incCol
	if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
		return false // out of board
	}
	dstLoc := dstRow*8 + dstCol
	return dstLoc == target
}
