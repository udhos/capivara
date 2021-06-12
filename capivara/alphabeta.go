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
	rootScore      []*boardScore // root children scores
}

func rootAlphaBeta(ab *alphaBetaState, b board, depth int, addChildren bool) (float32, move, string) {
	if depth < 1 {
		return relativeMaterial(b, addChildren), nullMove, "invalid-depth"
	}
	if b.otherKingInCheck() {
		return alphabetaMax, nullMove, "checkmate"
	}
	/*
		children := ab.children
		countChildren := b.generateChildren(children)
	*/
	countChildren := len(ab.rootScore)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin, nullMove, "checkmated" // checkmated
		}
		return 0, nullMove, "draw"
	}
	//firstChild := len(children.pool) - countChildren
	if countChildren == 1 {
		// in the root board, if there is a single possible move,
		// we can skip calculations and immediately return the move.
		// score is of course bogus in this case.
		ab.singleChildren = true
		return relativeMaterial(b, addChildren), ab.rootScore[0].b.lastMove, ""
	}

	var bestMove move
	var alpha float32 = alphabetaMin
	var beta float32 = alphabetaMax

	// handle first child
	{
		//child := children.pool[firstChild]
		child := ab.rootScore[0].b
		score := alphaBeta(ab, child, -beta, -alpha, depth-1, addChildren)
		score = -score
		ab.rootScore[0].score = score // update score for move ordering in next depth
		ab.nodes += int64(countChildren)
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
	//for _, child := range children.pool[firstChild+1:] {
	for _, rs := range ab.rootScore {
		child := rs.b
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
		rs.score = score // update score for move ordering in next depth
		ab.nodes += int64(countChildren)
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

func alphaBeta(ab *alphaBetaState, b board, alpha, beta float32, depth int, addChildren bool) float32 {

	children := ab.children

	if depth < 1 {
		return relativeMaterial(b, addChildren)
	}

	countChildren := b.generateChildren(children)
	if countChildren == 0 {
		if b.kingInCheck() {
			return alphabetaMin // checkmated
		}
		return 0 // draw
	}

	firstChild := len(children.pool) - countChildren
	lastChildren := children.pool[firstChild:]

	for _, child := range lastChildren {
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
		ab.nodes += int64(countChildren)
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
