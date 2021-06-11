package main

import "log"

var pieceSquareBlackKing = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 15, 15, 0, 0, 0,
	0, 0, 0, 15, 15, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 5, 15, 5, 5, 5, 15, 10,
	15, 10, 20, 10, 10, 10, 20, 15,
}

var pieceSquareBlackQueen = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pieceSquareBlackRook = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 0, 5, 0, 0,
}

var pieceSquareBlackBishop = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pieceSquareBlackKnight = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 10, 10, 10, 5, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pieceSquareBlackPawn = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	30, 30, 30, 30, 30, 30, 30, 30,
	10, 10, 10, 20, 20, 10, 10, 10,
	5, 5, 5, 15, 15, 5, 5, 5,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 5, 5, 0, 0, 5, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pieceSquareTableBlack = [6][64]int16{
	pieceSquareBlackKing,   // king
	pieceSquareBlackQueen,  // queen
	pieceSquareBlackRook,   // rook
	pieceSquareBlackBishop, // bishop
	pieceSquareBlackKnight, // knight
	pieceSquareBlackPawn,   // pawn
}

// pieceSquareTable: color => piece => location => value
var pieceSquareTable = [2][6][64]int16{
	pieceSquareTableBlack, // white - will mirror from black
	pieceSquareTableBlack, // black
}

func mirrorPieceSquareTable() {
	log.Printf("mirrorPieceSquareTable: mirroring white from black")
	for k := 0; k < 6; k++ {
		for row := 0; row < 8; row++ {
			for col := 0; col < 8; col++ {
				locWhite := row*8 + col
				locBlack := (7-row)*8 + col
				pieceSquareTable[colorWhite][k][locWhite] = pieceSquareTable[colorBlack][k][locBlack]
			}
		}
	}
	log.Printf("mirrorPieceSquareTable: done")
}
