package main

import (
	"fmt"
	"sort"
)

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

const (
	negamaxMin = -1000.0
	negamaxMax = 1000.0
)

type negamaxState struct {
	nodes int
}

func rootNegamax(nega *negamaxState, b board, depth int, path []string) (float32, string, []string) {
	if depth < 1 {
		return relativeMaterial(b), "invalid-depth", path
	}
	if b.otherKingInCheck() {
		return negamaxMax, "checkmate", nil
	}
	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return negamaxMin, "checkmated", path // checkmated
		}
		return 0, "draw", path
	}
	if len(children) == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		return relativeMaterial(children[0]), children[0].lastMove, path
	}

	var max float32 = negamaxMin
	/*
		var best string
		var bestPath []string
	*/

	var negaChildren []negaChild

	for _, child := range children {
		score, childPath := negamax(nega, child, depth-1, append(path, child.lastMove))
		score = -score
		nega.nodes += len(children)
		fmt.Printf("rootNegamax: depth=%d nodes=%d score=%v move: %s path: %s\n", depth, nega.nodes, score, child.lastMove, childPath)
		negaChildren = append(negaChildren, negaChild{b: child, score: score, path: childPath, nodes: len(children)})
		/*
			if score >= max {
				max = score
				best = child.lastMove
				bestPath = childPath
			}
		*/
	}

	fmt.Println()

	sort.Slice(negaChildren, func(i, j int) bool { return len(negaChildren[i].path) < len(negaChildren[j].path) })
	sort.SliceStable(negaChildren, func(i, j int) bool { return negaChildren[i].score > negaChildren[j].score })

	for _, c := range negaChildren {
		fmt.Printf("rootNegamax: depth=%d nodes=%d score=%v move: %s path: %s\n", depth, c.nodes, c.score, c.b.lastMove, c.path)
	}

	if negaChildren[0].score > max {
		max = negaChildren[0].score
	}

	return max, negaChildren[0].b.lastMove, negaChildren[0].path
}

type negaChild struct {
	b     board
	score float32
	path  []string
	nodes int
}

func negamax(nega *negamaxState, b board, depth int, path []string) (float32, []string) {
	if depth < 1 {
		return relativeMaterial(b), path
	}

	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return negamaxMin, path // checkmated
		}
		return 0, path // draw
	}

	var max float32 = negamaxMin
	var bestPath []string

	for _, child := range children {
		score, childPath := negamax(nega, child, depth-1, append(path, child.lastMove))
		score = -score
		nega.nodes += len(children)
		if score >= max {
			max = score
			bestPath = childPath
		}
	}
	return max, bestPath
}
