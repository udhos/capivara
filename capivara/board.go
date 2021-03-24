package main

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

	b.delPiece(i, j)

	loc := i*8 + j
	//fmt.Printf("addPiece: %dx%d=%d color=%d kind=%s\n", i, j, loc, p.color(), p.kindLetter())
	b.square[loc] = p

	// record king position
	if p.kind() == whiteKing {
		b.king[p.color()] = loc
	}

	b.materialValue[p.color()] += p.materialValue() // piece material value enters board
}

func (b *board) delPiece(i, j location) piece {
	loc := i*8 + j
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
