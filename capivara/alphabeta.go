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
}

func rootAlphaBeta(ab *alphaBetaState, b board, depth int, path []string, addChildren bool) (float32, string, []string) {
	if depth < 1 {
		return relativeMaterial(b, addChildren), "invalid-depth", path
	}
	if b.otherKingInCheck() {
		return alphabetaMax, "checkmate", nil
	}
	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return alphabetaMin, "checkmated", path // checkmated
		}
		return 0, "draw", path
	}
	if len(children) == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		ab.singleChildren = true
		return relativeMaterial(children[0], addChildren), children[0].lastMove, path
	}

	var bestPath []string
	var bestMove string
	var alpha float32 = alphabetaMin
	var beta float32 = alphabetaMax

	// handle first child
	{
		child := children[0]
		score, childPath := alphaBeta(ab, child, -beta, -alpha, depth-1, append(path, child.lastMove), addChildren)
		score = -score
		ab.nodes += int64(len(children))
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
	for _, child := range children[1:] {
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
		ab.nodes += int64(len(children))
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
	if depth < 1 {
		return relativeMaterial(b, addChildren), path
	}

	children := b.generateChildren([]board{})
	if len(children) == 0 {
		if b.kingInCheck() {
			return alphabetaMin, path // checkmated
		}
		return 0, path // draw
	}

	var bestPath []string

	for _, child := range children {
		if !ab.deadline.IsZero() {
			// there is a timer
			if ab.deadline.Before(time.Now()) {
				// timer has expired
				ab.cancelled = true
				return 0, nil
			}
		}
		score, childPath := alphaBeta(ab, child, -beta, -alpha, depth-1, append(path, child.lastMove), addChildren)
		score = -score
		ab.nodes += int64(len(children))
		if score >= beta {
			return beta, childPath
		}
		if score > alpha {
			alpha = score
			bestPath = childPath
		}
	}

	return alpha, bestPath
}
