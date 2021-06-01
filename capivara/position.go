package main

var defaultPositionWeight = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0, // rank 1
	0, 5, 5, 5, 5, 5, 5, 0, // rank 2
	0, 5, 10, 10, 10, 10, 5, 0, // rank 3
	0, 5, 10, 15, 15, 10, 5, 0, // rank 4
	0, 5, 10, 15, 15, 10, 5, 0, // rank 5
	0, 5, 10, 10, 10, 10, 5, 0, // rank 6
	0, 5, 5, 5, 5, 5, 5, 0, // rank 7
	0, 0, 0, 0, 0, 0, 0, 0, // rank 8
}

var positionTable = [6][64]int16{
	defaultPositionWeight, // king
	defaultPositionWeight, // queen
	defaultPositionWeight, // rook
	defaultPositionWeight, // bishop
	defaultPositionWeight, // knight
	defaultPositionWeight, // pawn
}
