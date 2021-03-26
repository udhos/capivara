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

const negamaxMin = -1000.0

type negamaxState struct {
	nodes int
}

func rootNegamax(nega *negamaxState, b board, depth int, path []string) (float32, string, []string) {
	if depth < 1 {
		return relativeMaterial(b), "invalid-depth", path
	}
	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return negamaxMin * float32(colorToSignal(b.turn)), "checkmated", path
		}
		return 0, "draw", path
	}

	var max float32 = negamaxMin
	var best string
	var bestPath []string

	for _, child := range children {
		score, childPath := negamax(nega, child, depth-1, append(path, child.lastMove))
		score = -score
		nega.nodes += len(children)
		fmt.Printf("rootNegamax: depth=%d nodes=%d score=%v move: %s path: %s\n", depth, nega.nodes, score, child.lastMove, childPath)
		if score > max {
			max = score
			best = child.lastMove
			bestPath = childPath
		}
	}
	return max, best, bestPath
}

func negamax(nega *negamaxState, b board, depth int, path []string) (float32, []string) {
	if depth < 1 {
		return relativeMaterial(b), path
	}

	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return negamaxMin * float32(colorToSignal(b.turn)), path // checkmated
		}
		return relativeMaterial(b), path // draw
	}

	var max float32 = negamaxMin
	var bestPath []string

	for _, child := range children {
		score, childPath := negamax(nega, child, depth-1, append(path, child.lastMove))
		score = -score
		nega.nodes += len(children)
		//fmt.Printf("negamax: depth=%d score=%v move: %s\n", depth, score, child.lastMove)
		if score > max {
			max = score
			bestPath = childPath
		}
	}
	return max, bestPath
}
