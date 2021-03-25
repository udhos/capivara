package main

import "fmt"

func rootNegamax(b board, depth int) (float32, string) {
	if depth == 0 {
		return b.getMaterialValue(), "move?"
	}
	children := b.generateChildren([]board{})
	if len(children) == 0 {
		return b.getMaterialValue(), "move?"
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
		return b.getMaterialValue()
	}

	children := b.generateChildren([]board{})
	if len(children) == 0 {
		return b.getMaterialValue()
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
