package main

func (b *board) otherKingInCheck() bool {
	b.turn = colorInverse(b.turn)
	check := b.kingInCheck()
	b.turn = colorInverse(b.turn)
	return check
}

func (b *board) kingInCheck() bool {
	return b.anyPieceAttacks(b.king[b.turn])
}

func (b *board) anyPieceAttacks(loc location) bool {

	// pawn

	if b.findAttackFromPawn(loc) {
		return true
	}

	// king

	if b.findAttackFromKing(loc) {
		return true
	}

	// knight

	if b.findAttackFromKnight(loc) {
		return true
	}

	// rook or queen

	if b.findAttackFromHV(loc) {
		return true
	}

	// bishop or queen

	if b.findAttackFromDiagonal(loc) {
		return true
	}

	return false
}

func (b *board) findAttackFromPawn(kingLoc location) bool {
	row := kingLoc / 8
	col := kingLoc % 8
	signal := colorToSignal(b.turn) // 0=>1 1=>-1
	srcRow := int(row) + signal
	if srcRow < 1 || srcRow > 6 {
		return false
	}
	if col > 0 {
		//log.Printf("king %c%c pawn %c%c", col+'a', row+'1', (col-1)+'a', srcRow+'1')
		srcLoc := srcRow*8 + int(col) - 1
		p := b.square[srcLoc]
		if p != pieceNone && p.color() != b.turn && p.kind() == whitePawn {
			return true
		}
	}
	if col < 7 {
		//log.Printf("king %c%c pawn %c%c", col+'a', row+'1', (col+1)+'a', srcRow+'1')
		srcLoc := srcRow*8 + int(col) + 1
		p := b.square[srcLoc]
		if p != pieceNone && p.color() != b.turn && p.kind() == whitePawn {
			return true
		}
	}
	return false
}

func (b *board) findAttackFromKing(trg location) bool {
	trgRow := trg / 8
	trgCol := trg % 8

	if trgRow > 0 {
		// row-1

		// row-1 col
		srcRow := trgRow - 1
		if p := b.square[srcRow*8+trgCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
			return true
		}

		if trgCol > 0 {
			// row-1 col-1
			srcCol := trgCol - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
				return true
			}
		}

		if trgCol < 7 {
			// row-1 col+1
			srcCol := trgCol + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
				return true
			}
		}
	}

	if trgRow < 7 {
		// row+1

		// row+1 col
		srcRow := trgRow + 1
		if p := b.square[srcRow*8+trgCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
			return true
		}

		if trgCol > 0 {
			// row+1 col-1
			srcCol := trgCol - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
				return true
			}
		}

		if trgCol < 7 {
			// row+1 col+1
			srcCol := trgCol + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
				return true
			}
		}
	}

	if trgCol > 0 {
		// row col-1
		srcCol := trgCol - 1
		if p := b.square[trgRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
			return true
		}
	}

	if trgCol < 7 {
		// row col+1
		srcCol := trgCol + 1
		if p := b.square[trgRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKing {
			return true
		}
	}

	return false
}

func (b *board) findAttackFromKnight(kingLoc location) bool {
	trgRow := kingLoc / 8
	trgCol := kingLoc % 8

	if trgRow > 1 {
		// row-2
		srcRow := trgRow - 2

		if trgCol > 0 {
			// row-2 col-1
			srcCol := trgCol - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}

		if trgCol < 7 {
			// row-2 col+1
			srcCol := trgCol + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}
	}

	if trgRow < 6 {
		// row+2
		srcRow := trgRow + 2

		if trgCol > 0 {
			// row+2 col-1
			srcCol := trgCol - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}

		if trgCol < 7 {
			// row+2 col+1
			srcCol := trgCol + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}
	}

	if trgCol > 1 {
		// col-2
		srcCol := trgCol - 2

		if trgRow > 0 {
			// row-1 col-2
			srcRow := trgRow - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}

		if trgRow < 7 {
			// row+1 col-2
			srcRow := trgRow + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}
	}

	if trgCol < 6 {
		// col+2
		srcCol := trgCol + 2

		if trgRow > 0 {
			// row-1 col+2
			srcRow := trgRow - 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}

		if trgRow < 7 {
			// row+1 col+2
			srcRow := trgRow + 1
			if p := b.square[srcRow*8+srcCol]; p != pieceNone && p.color() != b.turn && p.kind() == whiteKnight {
				return true
			}
		}
	}

	return false
}

func (b *board) findAttackFromHV(kingLoc location) bool {
	kingRow := kingLoc / 8
	kingCol := kingLoc % 8

	// up
	for row := kingRow; row < 7; {
		row++
		loc := row*8 + kingCol
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteRook || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// down
	for row := kingRow; row > 0; {
		row--
		loc := row*8 + kingCol
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteRook || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// left
	for col := kingCol; col > 0; {
		col--
		loc := kingRow*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteRook || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// right
	for col := kingCol; col < 7; {
		col++
		loc := kingRow*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteRook || kind == whiteQueen {
				return true
			}
		}
		break
	}

	return false
}

func (b *board) findAttackFromDiagonal(kingLoc location) bool {
	kingRow := kingLoc / 8
	kingCol := kingLoc % 8

	// - -
	for row, col := kingRow, kingCol; row > 0 && col > 0; {
		row--
		col--
		loc := row*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteBishop || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// - +
	for row, col := kingRow, kingCol; row > 0 && col < 7; {
		row--
		col++
		loc := row*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteBishop || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// + -
	for row, col := kingRow, kingCol; row < 7 && col > 0; {
		row++
		col--
		loc := row*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteBishop || kind == whiteQueen {
				return true
			}
		}
		break
	}

	// + +
	for row, col := kingRow, kingCol; row < 7 && col < 7; {
		row++
		col++
		loc := row*8 + col
		p := b.square[loc]
		if p == pieceNone {
			continue
		}
		if p.color() != b.turn {
			// other color piece
			if kind := p.kind(); kind == whiteBishop || kind == whiteQueen {
				return true
			}
		}
		break
	}

	return false
}
