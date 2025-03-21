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
func relativeMaterial(children *boardPool, b board, addChildren bool) float32 {
	relative := float32(colorToSignal(b.turn)) * b.getMaterialValue()
	if addChildren {
		countChildren := b.generateChildren(children)
		relative += float32(countChildren) / 100.0
		children.drop(countChildren)
	}
	return relative
}

const (
	negamaxMin = -1000.0
	negamaxMax = 1000.0
)

type negamaxState struct {
	nodes      int64
	children   *boardPool
	showSearch bool
}

func rootNegamax(nega *negamaxState, b board, depth int, addChildren bool) (float32, move, string) {
	if depth < 1 {
		return relativeMaterial(nega.children, b, addChildren), nullMove, "invalid-depth"
	}
	if b.otherKingInCheck() {
		return negamaxMax, nullMove, "checkmate"
	}

	children := nega.children
	countChildren := b.generateChildren(children)

	if countChildren == 0 {
		if b.kingInCheck() {
			return negamaxMin, nullMove, "checkmated" // checkmated
		}
		return 0, nullMove, "draw"
	}

	nega.nodes += int64(countChildren)

	firstChild := len(children.pool) - countChildren
	if countChildren == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		return relativeMaterial(children, children.pool[firstChild], addChildren), children.pool[firstChild].lastMove, ""
	}

	var maxScore float32 = negamaxMin

	var negaChildren []negaChild

	lastChildren := children.pool[firstChild:]

	for _, child := range lastChildren {
		score := negamax(nega, child, depth-1, addChildren)
		score = -score
		if nega.showSearch {
			fmt.Printf("rootNegamax: depth=%d nodes=%d score=%v move: %s\n", depth, nega.nodes, score, child.lastMove)
		}
		negaChildren = append(negaChildren, negaChild{b: child, score: score, nodes: countChildren})
	}

	sort.SliceStable(negaChildren, func(i, j int) bool { return negaChildren[i].score > negaChildren[j].score })

	if negaChildren[0].score > maxScore {
		maxScore = negaChildren[0].score
	}

	return maxScore, negaChildren[0].b.lastMove, ""
}

type negaChild struct {
	b     board
	score float32
	nodes int
}

func negamax(nega *negamaxState, b board, depth int, addChildren bool) float32 {

	children := nega.children

	if depth < 1 {
		return relativeMaterial(children, b, addChildren)
	}

	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return negamaxMin // checkmated
		}
		return 0 // draw
	}

	nega.nodes += int64(countChildren)

	var maxScore float32 = negamaxMin

	firstChild := len(children.pool) - countChildren
	lastChildren := children.pool[firstChild:]

	for _, child := range lastChildren {
		score := negamax(nega, child, depth-1, addChildren)
		score = -score
		if score >= maxScore {
			maxScore = score
		}
	}

	children.drop(countChildren)
	return maxScore
}
