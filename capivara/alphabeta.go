package main

import (
	"fmt"
	"time"
)

const (
	alphabetaMin = -1000.0
	alphabetaMax = 1000.0
)

type alphaBetaState struct {
	nodes          int64
	showSearch     bool
	deadline       time.Time
	cancelled      bool
	singleChildren bool
	children       *boardPool
}

func rootAlphaBeta(ab *alphaBetaState, b board, depth int, path []string, addChildren bool) (float32, string, []string) {
	if depth < 1 {
		return relativeMaterial(ab.children, b, addChildren), "invalid-depth", path
	}
	if b.otherKingInCheck() {
		return alphabetaMax, "checkmate", nil
	}
	children := ab.children
	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin, "checkmated", path // checkmated
		}
		return 0, "draw", path
	}
	if countChildren == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		ab.singleChildren = true
		return relativeMaterial(children, b, addChildren), ab.children.pool[0].lastMove, path
	}

	var bestPath []string
	var bestMove string
	var alpha float32 = alphabetaMin
	var beta float32 = alphabetaMax

	// handle first child
	{
		child := children.pool[0]
		score, childPath := alphaBeta(ab, child, -beta, -alpha, depth-1, append(path, child.lastMove), addChildren)
		score = -score
<<<<<<< HEAD
		ab.nodes += countChildren
=======
		ab.nodes += int64(len(children))
>>>>>>> main
		if ab.showSearch {
			fmt.Printf("rootAlphaBeta: depth=%d nodes=%d score=%v move: %s path: %s\n", depth, ab.nodes, score, child.lastMove, childPath)
		}
		if score >= beta {
			return beta, child.lastMove, childPath
		}

		// pick first child
		alpha = score
		bestMove = child.lastMove
		bestPath = childPath
	}

	// scan remaining children
	for _, child := range children.pool[1:] {
		if !ab.deadline.IsZero() {
			// there is a timer
			if ab.deadline.Before(time.Now()) {
				// timer has expired
				ab.cancelled = true
				return 0, "", nil
			}
		}
		score, childPath := alphaBeta(ab, child, -beta, -alpha, depth-1, append(path, child.lastMove), addChildren)
		score = -score
<<<<<<< HEAD
		ab.nodes += countChildren
=======
		ab.nodes += int64(len(children))
>>>>>>> main
		if ab.showSearch {
			fmt.Printf("rootAlphaBeta: depth=%d nodes=%d score=%v move: %s path: %s\n", depth, ab.nodes, score, child.lastMove, childPath)
		}
		if score >= beta {
			return beta, child.lastMove, childPath
		}
		if score > alpha {
			alpha = score
			bestMove = child.lastMove
			bestPath = childPath
		}
	}

	return alpha, bestMove, bestPath
}

func alphaBeta(ab *alphaBetaState, b board, alpha, beta float32, depth int, path []string, addChildren bool) (float32, []string) {

	children := ab.children

	if depth < 1 {
		return relativeMaterial(children, b, addChildren), path
	}

	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin, path // checkmated
		}
		return 0, path // draw
	}

	lastChildren := children.pool[len(children.pool)-countChildren:]

	var bestPath []string

	for _, child := range lastChildren {
		if !ab.deadline.IsZero() {
			// there is a timer
			if ab.deadline.Before(time.Now()) {
				// timer has expired
				ab.cancelled = true
				children.drop(countChildren)
				return 0, nil
			}
		}
		score, childPath := alphaBeta(ab, child, -beta, -alpha, depth-1, append(path, child.lastMove), addChildren)
		score = -score
		ab.nodes += int64(countChildren)
		if score >= beta {
			children.drop(countChildren)
			return beta, childPath
		}
		if score > alpha {
			alpha = score
			bestPath = childPath
		}
	}

	children.drop(countChildren)
	return alpha, bestPath
}
