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
	lastMove      move
	zobristValue  zobristKey
	parent        *board
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
	//w := positionWeight[loc] * int16(colorToSignal(p.color()))
	value := p.materialValue(loc)
	b.materialValue[p.color()] += value // piece material value enters board
	//log.Printf("add: loc=%d material=%d board=%d", loc, value, b.materialValue[p.color()])
	if p.kind() == whiteKing {
		// record king new position
		b.king[p.color()] = loc
	}

	b.zobristUpdatePiece(int(loc), p) // add zobrist value after adding piece
}

func (b *board) delPieceLoc(loc location) piece {
	p := b.square[loc]
	if p.kind() != pieceNone {
		//w := positionWeight[loc] * int16(colorToSignal(p.color()))
		value := p.materialValue(loc)
		b.materialValue[p.color()] -= value // piece material value leaves board
		//log.Printf("del: loc=%d material=%d board=%d", loc, value, b.materialValue[p.color()])

		b.zobristUpdatePiece(int(loc), p) // remove zobrist value before removing piece

		b.square[loc] = pieceNone
	}
	return p
}

func (b board) getMaterialValue() float32 {
	wh := float32(b.materialValue[0])
	bl := float32(b.materialValue[1])
	return (wh + bl) / 100
}

func (b *board) generatePassantCapture(attackerLoc, targetLoc location, children *boardPool) int {
	attackerP := b.square[attackerLoc]
	attackerColor := attackerP.color()
	targetColor := b.square[targetLoc].color()

	if attackerP.kind() == whitePawn && attackerColor != targetColor {

		attackerSignal := colorToSignal(attackerColor)
		attackerDstLoc := targetLoc + location(8*attackerSignal)

		c, _ := b.createChild(children, attackerLoc, attackerDstLoc)

		c.delPieceLoc(targetLoc) // captured passant pawn

		return keepIfValid(children, c)
	}

	return 0
}

func (b board) generateChildren(children *boardPool) int {

	var countChildren int

	// generate en passant captures

	lastMove := b.lastMove
	if !lastMove.isNull() {
		step := lastMove.rankDelta()
		if step == 2 {
			trgDstLoc := lastMove.dst
			trgKind := b.square[trgDstLoc].kind()

			if trgKind == whitePawn {
				// it is pawn

				trgDstCol := trgDstLoc % 8
				trgDstRow := trgDstLoc / 8
				trgDstRow8 := 8 * trgDstRow

				if trgDstCol > 0 {
					// might be captured from left
					attackerLoc := location(trgDstCol - 1 + trgDstRow8)
					countChildren += b.generatePassantCapture(attackerLoc, trgDstLoc, children)
				}
				if trgDstCol < 7 {
					// might be captured from right
					attackerLoc := location(trgDstCol + 1 + trgDstRow8)
					countChildren += b.generatePassantCapture(attackerLoc, trgDstLoc, children)
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
		countChildren += b.generateChildrenPiece(children, loc, p)
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
				countChildren += b.generateCastlingLeft(children)
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
				countChildren += b.generateCastlingRight(children)
			}
		}
	}

	return countChildren
}

func (b board) generateCastlingLeft(children *boardPool) int {

	row := 7 * location(b.turn)
	row8 := 8 * row
	kingSrc := row8 + 4 // E
	kingDst := row8 + 2 // C
	rookSrc := row8     // A
	rookDst := row8 + 3 // D

	//child := b // copy board

	// move
	b.square[kingDst] = b.square[kingSrc]
	b.square[rookDst] = b.square[rookSrc]
	b.square[kingSrc] = pieceNone
	b.square[rookSrc] = pieceNone

	// record king new position
	b.king[b.turn] = kingDst

	// disable castling
	b.zobristUpdateCastling()
	b.flags[b.turn] |= lostCastlingLeft | lostCastlingRight
	b.zobristUpdateCastling()

	b.zobristUpdateTurn()
	b.turn = colorInverse(b.turn) // switch color
	b.zobristUpdateTurn()

	//b.lastMove = moveToStr(kingSrc, kingDst, pieceNone) // record last move
	b.zobristUpdateEnPassant()
	b.lastMove = move{src: kingSrc, dst: kingDst} // record last move
	b.zobristUpdateEnPassant()

	//return b.recordIfValid(children, child)
	// no need to verify king in check since castling conditions
	// previously required king target square is free from attack
	children.push(&b)
	return 1
}

func (b board) generateCastlingRight(children *boardPool) int {
	row := 7 * location(b.turn)
	row8 := 8 * row
	kingSrc := row8 + 4 // E
	kingDst := row8 + 6 // G
	rookSrc := row8 + 7 // H
	rookDst := row8 + 5 // F

	//child := b // copy board

	// move
	b.square[kingDst] = b.square[kingSrc]
	b.square[rookDst] = b.square[rookSrc]
	b.square[kingSrc] = pieceNone
	b.square[rookSrc] = pieceNone

	// record king new position
	b.king[b.turn] = kingDst

	// disable castling
	b.zobristUpdateCastling()
	b.flags[b.turn] |= lostCastlingLeft | lostCastlingRight
	b.zobristUpdateCastling()

	b.zobristUpdateTurn()
	b.turn = colorInverse(b.turn) // switch color
	b.zobristUpdateTurn()

	//b.lastMove = moveToStr(kingSrc, kingDst, pieceNone) // record last move
	b.zobristUpdateEnPassant()
	b.lastMove = move{src: kingSrc, dst: kingDst} // record last move
	b.zobristUpdateEnPassant()

	//return b.recordIfValid(children, child)
	// no need to verify king in check since castling conditions
	// previously required king target square is free from attack
	children.push(&b)
	return 1
}

func (b board) generateChildrenPiece(children *boardPool, loc location, p piece) int {
	var countChildren int

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
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					countChildren += b.recordMoveIfValid(children, loc, location(dstLoc))
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
				countChildren += b.recordMoveIfValid(children, loc, location(dstRowLoc))
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
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					countChildren += b.recordMoveIfValid(children, loc, location(dstLoc))
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
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteQueen)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteRook)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteBishop)
					countChildren += b.recordPromotionIfValid(children, loc, location(dstLoc), piece(color<<3)+whiteKnight)
				} else {
					countChildren += b.recordMoveIfValid(children, loc, location(dstLoc))
				}
			}
		}

	case whiteQueen: // white + black
		countChildren += b.generateSliding(children, loc, 0, 1)
		countChildren += b.generateSliding(children, loc, 1, 1)
		countChildren += b.generateSliding(children, loc, 1, 0)
		countChildren += b.generateSliding(children, loc, 1, -1)
		countChildren += b.generateSliding(children, loc, 0, -1)
		countChildren += b.generateSliding(children, loc, -1, -1)
		countChildren += b.generateSliding(children, loc, -1, 0)
		countChildren += b.generateSliding(children, loc, -1, 1)

	case whiteRook: // white + black
		countChildren += b.generateSlidingRook(children, loc, 0, 1)
		countChildren += b.generateSlidingRook(children, loc, 1, 0)
		countChildren += b.generateSlidingRook(children, loc, 0, -1)
		countChildren += b.generateSlidingRook(children, loc, -1, 0)

	case whiteBishop: // white + black
		countChildren += b.generateSliding(children, loc, 1, 1)
		countChildren += b.generateSliding(children, loc, 1, -1)
		countChildren += b.generateSliding(children, loc, -1, -1)
		countChildren += b.generateSliding(children, loc, -1, 1)

	case whiteKing: // white + black
		countChildren += b.generateRelativeKing(children, loc, 0, 1)
		countChildren += b.generateRelativeKing(children, loc, 1, 1)
		countChildren += b.generateRelativeKing(children, loc, 1, 0)
		countChildren += b.generateRelativeKing(children, loc, 1, -1)
		countChildren += b.generateRelativeKing(children, loc, 0, -1)
		countChildren += b.generateRelativeKing(children, loc, -1, -1)
		countChildren += b.generateRelativeKing(children, loc, -1, 0)
		countChildren += b.generateRelativeKing(children, loc, -1, 1)

	case whiteKnight: // white + black
		countChildren += b.generateRelative(children, loc, -1, 2)
		countChildren += b.generateRelative(children, loc, 1, 2)
		countChildren += b.generateRelative(children, loc, 2, -1)
		countChildren += b.generateRelative(children, loc, 2, 1)
		countChildren += b.generateRelative(children, loc, -1, -2)
		countChildren += b.generateRelative(children, loc, 1, -2)
		countChildren += b.generateRelative(children, loc, -2, -1)
		countChildren += b.generateRelative(children, loc, -2, 1)
	}

	return countChildren
}

func (b board) recordIfValid(children *boardPool, child board) int {
	if child.otherKingInCheck() {
		return 0 // drop invalid move 'child'
	}
	children.push(&child) // record
	return 1
}

func keepIfValid(children *boardPool, child *board) int {
	if child.otherKingInCheck() {
		children.drop(1)
		return 0 // drop invalid move 'child'
	}
	return 1 // keep
}

func (b board) generateSliding(children *boardPool, src, incRow, incCol location) int {
	var countChildren int

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
			countChildren += b.recordMoveIfValid(children, src, dstLoc)
			continue
		}
		if dstP.color() != b.turn {
			// capture opponent piece
			countChildren += b.recordMoveIfValid(children, src, dstLoc)
		}
		break
	}

	return countChildren
}

func (b board) generateSlidingRook(children *boardPool, src, incRow, incCol location) int {
	var countChildren int

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
			countChildren += b.recordMoveIfValidRook(children, src, dstLoc)
			continue
		}
		if dstP.color() != b.turn {
			// capture opponent piece
			countChildren += b.recordMoveIfValidRook(children, src, dstLoc)
		}
		break
	}

	return countChildren
}

func (b board) generateRelative(children *boardPool, src, incRow, incCol location) int {
	dstRow := src / 8
	dstCol := src % 8

	dstRow += incRow
	dstCol += incCol
	if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
		return 0 // out of board
	}
	dstLoc := dstRow*8 + dstCol
	dstP := b.square[dstLoc]
	if dstP == pieceNone {
		// empty square
		return b.recordMoveIfValid(children, src, dstLoc)
	}
	if dstP.color() != b.turn {
		// capture opponent piece
		return b.recordMoveIfValid(children, src, dstLoc)
	}

	// blocked by same color piece

	return 0
}

func (b board) generateRelativeKing(children *boardPool, src, incRow, incCol location) int {
	dstRow := src / 8
	dstCol := src % 8

	dstRow += incRow
	dstCol += incCol
	if dstRow < 0 || dstRow > 7 || dstCol < 0 || dstCol > 7 {
		return 0 // out of board
	}
	dstLoc := dstRow*8 + dstCol
	dstP := b.square[dstLoc]
	if dstP == pieceNone {
		// empty square
		return b.recordMoveIfValidKing(children, src, dstLoc)
	}
	if dstP.color() != b.turn {
		// capture opponent piece
		return b.recordMoveIfValidKing(children, src, dstLoc)
	}

	// blocked by same color piece

	return 0
}

func (b board) newChild(src, dst location) (board, piece) {
	//child := b                                      // copy board
	p := b.delPieceLoc(src) // take piece from board
	b.addPieceLoc(dst, p)   // put piece on board

	b.zobristUpdateTurn()
	b.turn = colorInverse(b.turn) // switch color
	b.zobristUpdateTurn()

	b.zobristUpdateEnPassant()
	b.lastMove = move{src: src, dst: dst} // record move
	b.zobristUpdateEnPassant()

	return b, p
}

func (b *board) createChild(children *boardPool, src, dst location) (*board, piece) {

	children.push(b)         // copy
	child := children.last() // get address
	child.parent = b         // point to parent

	p := child.delPieceLoc(src) // take piece from board
	child.addPieceLoc(dst, p)   // put piece on board

	child.zobristUpdateTurn()
	child.turn = colorInverse(b.turn) // switch color
	child.zobristUpdateTurn()

	child.zobristUpdateEnPassant()
	child.lastMove = move{src: src, dst: dst} // record move
	child.zobristUpdateEnPassant()

	return child, p
}

func (b board) recordMoveIfValid(children *boardPool, src, dst location) int {
	child, _ := b.newChild(src, dst)
	return b.recordIfValid(children, child)
}

func (b board) recordMoveIfValidKing(children *boardPool, src, dst location) int {
	child, p := b.newChild(src, dst)

	// king moved, then disable castling
	child.flags[p.color()] |= lostCastlingLeft | lostCastlingRight

	return b.recordIfValid(children, child)
}

func (b board) recordMoveIfValidRook(children *boardPool, src, dst location) int {
	child, p := b.newChild(src, dst)

	// rook moved, then disable castling
	firstRow := 7 * location(p.color()) // 0=>0 1=>7
	row := src / 8
	if row == firstRow {
		// rook in initial row
		col := src % 8
		switch col {
		case 0: // left rook
			b.zobristUpdateCastling()
			child.flags[p.color()] |= lostCastlingLeft
			b.zobristUpdateCastling()
		case 7: // right rook
			b.zobristUpdateCastling()
			child.flags[p.color()] |= lostCastlingRight
			b.zobristUpdateCastling()
		}
	}

	return b.recordIfValid(children, child)
}

func (b board) recordPromotionIfValid(children *boardPool, src, dst location, p piece) int {
	//child := b                              // copy board
	b.delPieceLoc(src)    // take pawn from board
	b.addPieceLoc(dst, p) // put new piece on board

	b.zobristUpdateTurn()
	b.turn = colorInverse(b.turn) // switch color
	b.zobristUpdateTurn()

	//b.lastMove = moveToStr(src, dst, p) // record move
	b.zobristUpdateEnPassant()
	b.lastMove = move{src: src, dst: dst, promotion: p} // record move
	b.zobristUpdateEnPassant()

	return b.recordIfValid(children, b)
}
