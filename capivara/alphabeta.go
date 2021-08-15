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

func rootAlphaBeta(ab *alphaBetaState, b *board, depth int, addChildren bool) (float32, move, string) {
	if depth < 1 {
		if b.lastMove.isQuiescent() {
			return relativeMaterial(ab.children, b, addChildren), nullMove, "invalid-depth"
		}
	}
	if b.otherKingInCheck() {
		return alphabetaMax, nullMove, "checkmate"
	}
	children := ab.children
	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin, nullMove, "checkmated" // checkmated
		}
		return 0, nullMove, "draw"
	}

	ab.nodes += int64(countChildren)

	firstChild := len(children.pool) - countChildren
	if countChildren == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		ab.singleChildren = true
		return relativeMaterial(children, b, addChildren), ab.children.pool[firstChild].lastMove, ""
	}

	var bestMove move
	var alpha float32 = alphabetaMin
	var beta float32 = alphabetaMax

	// handle first child
	{
		child := &children.pool[firstChild]
		score := alphaBeta(ab, child, -beta, -alpha, depth-1, addChildren)
		score = -score
		if ab.showSearch {
			fmt.Printf("rootAlphaBeta: depth=%d nodes=%d score=%v move: %s\n", depth, ab.nodes, score, child.lastMove)
		}
		if score >= beta {
			return beta, child.lastMove, ""
		}

		// pick first child
		alpha = score
		bestMove = child.lastMove
	}

	// scan remaining children
	for i := firstChild + 1; i < len(children.pool); i++ {
		child := &children.pool[i]
		if !ab.deadline.IsZero() {
			// there is a timer
			if ab.deadline.Before(time.Now()) {
				// timer has expired
				ab.cancelled = true
				return 0, nullMove, ""
			}
		}
		score := alphaBeta(ab, child, -beta, -alpha, depth-1, addChildren)
		score = -score
		if ab.showSearch {
			fmt.Printf("rootAlphaBeta: depth=%d nodes=%d score=%v move: %s\n", depth, ab.nodes, score, child.lastMove)
		}
		if score >= beta {
			return beta, child.lastMove, ""
		}
		if score > alpha {
			alpha = score
			bestMove = child.lastMove
		}
	}

	return alpha, bestMove, ""
}

func alphaBeta(ab *alphaBetaState, b *board, alpha, beta float32, depth int, addChildren bool) float32 {

	children := ab.children

	if depth < 1 {
		if b.lastMove.isQuiescent() {
			return relativeMaterial(children, b, addChildren)
		}
	}

	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin // checkmated
		}
		return 0 // draw
	}

	ab.nodes += int64(countChildren)

	firstChild := len(children.pool) - countChildren

	for i := firstChild; i < len(children.pool); i++ {
		child := &children.pool[i]
		if !ab.deadline.IsZero() {
			// there is a timer
			if ab.deadline.Before(time.Now()) {
				// timer has expired
				ab.cancelled = true
				children.drop(countChildren)
				return 0
			}
		}
		score := alphaBeta(ab, child, -beta, -alpha, depth-1, addChildren)
		score = -score
		if score >= beta {
			children.drop(countChildren)
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	children.drop(countChildren)
	return alpha
}
