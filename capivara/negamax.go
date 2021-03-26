package main

import "fmt"

// negamax needs a relative material score.
//
// board.getMaterialValue() computes an absolute material score:
// the higher the better for the white player
//
// relativeMaterial(board) converts absolute material score to relative:
// the higher the better for the current player
func relativeMaterial(b board) float32 {
	return float32(colorToSignal(b.turn)) * b.getMaterialValue()
}

func rootNegamax(b board, depth int) (float32, string) {
	if depth < 1 {
		return relativeMaterial(b), "invalid-depth"
	}
	children := b.generateChildren([]board{})
	if len(children) == 0 {
		return relativeMaterial(b), "no-valid-move"
	}

	var max float32 = -1000.0
	var best string

	for _, child := range children {
		score := -negamax(child, depth-1)
		fmt.Printf("rootNegamax: depth=%d score=%v move: %s\n", depth, score, child.lastMove)
		if score > max {
			max = score
			best = child.lastMove
		}
	}
	return max, best
}

func negamax(b board, depth int) float32 {
	if depth == 0 {
		return relativeMaterial(b)
	}

	children := b.generateChildren([]board{})
	if len(children) == 0 {
		return relativeMaterial(b)
	}

	var max float32 = -1000.0

	for _, child := range children {
		score := -negamax(child, depth-1)
		//fmt.Printf("negamax: depth=%d score=%v move: %s\n", depth, score, child.lastMove)
		if score > max {
			max = score
		}
	}
	return max
}
