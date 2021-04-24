package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type command struct {
	name        string
	call        func(cmds []command, game *gameState, tokens []string)
	description string
}

var tableCmd = []command{
	{"ab", cmdAlphaBeta, "ab [depth] - alpha-beta search"},
	{"clear", cmdClear, "erase board"},
	{"help", cmdHelp, "show help"},
	{"load", cmdLoad, "load board from file"},
	{"move", cmdMove, "change piece position"},
	{"negamax", cmdNegamax, "negamax [depth] - negamax search"},
	{"play", cmdPlay, "play move"},
	{"perft", cmdPerft, "perft depth - count moves to depth"},
	{"reset", cmdReset, "reset board to initial position"},
	{"search", cmdSearch, "search [ms] - search"},
	{"switch", cmdSwitch, "switch turn"},
	{"undo", cmdUndo, "undo last played move"},
}

func cmdClear(cmds []command, game *gameState, tokens []string) {
	*game = newGame()
}

func cmdHelp(cmds []command, game *gameState, tokens []string) {
	fmt.Println("available commands:")
	for _, cmd := range cmds {
		fmt.Printf(" %s - %s\n", cmd.name, cmd.description)
	}
}

func cmdLoad(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: load filename\n")
		return
	}
	game.loadFromFile(tokens[1])
}

func cmdMove(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: move fromto\n")
		return
	}
	move := tokens[1]
	if len(move) < 4 || len(move) > 5 {
		fmt.Printf("usage: bad move length=%d: '%s'\n", len(move), move)
		return
	}
	from := strings.ToLower(move[:2])
	to := strings.ToLower(move[2:4])

	fmt.Printf("[%s][%s]\n", from, to)

	// from
	if len(from) != 2 {
		fmt.Printf("bad source format: [%s]\n", from)
		return
	}
	if from[0] < 'a' || from[0] > 'h' {
		fmt.Printf("bad source column letter: [%s]\n", from)
	}
	if from[1] < '1' || from[1] > '8' {
		fmt.Printf("bad source row digit: [%s]\n", from)
	}

	// to
	if len(to) != 2 {
		fmt.Printf("bad target format: [%s]\n", to)
		return
	}
	if to[0] < 'a' || to[0] > 'h' {
		fmt.Printf("bad target column letter: [%s]\n", to)
	}
	if to[1] < '1' || to[1] > '8' {
		fmt.Printf("bad target row digit: [%s]\n", to)
	}

	b := &game.history[len(game.history)-1]                       // will update in-place
	p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board

	if len(move) > 4 {
		// promotion
		promotion := move[4]
		kind := pieceKindFromLetter(rune(promotion))
		if kind != pieceNone {
			p = piece(b.turn<<3) + kind
		}
	}

	b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p) // put piece on board
}

func cmdNegamax(cmds []command, game *gameState, tokens []string) {
	depth := 4
	if len(tokens) > 1 {
		d, errConv := strconv.Atoi(tokens[1])
		if errConv == nil {
			depth = d
		}
	}
	fmt.Printf("negamax depth=%d\n", depth)
	last := len(game.history) - 1
	b := game.history[last]
	nega := negamaxState{}
	score, move, path := rootNegamax(&nega, b, depth, make([]string, 0, 20), game.addChildren)
	fmt.Printf("negamax: nodes=%d best score=%v move: %s path: %s\n", nega.nodes, score, move, path)
}

func cmdAlphaBeta(cmds []command, game *gameState, tokens []string) {
	depth := 4
	if len(tokens) > 1 {
		d, errConv := strconv.Atoi(tokens[1])
		if errConv == nil {
			depth = d
		}
	}
	fmt.Printf("alphabeta depth=%d\n", depth)
	last := len(game.history) - 1
	b := game.history[last]
	ab := alphaBetaState{showSearch: true}
	score, move, path := rootAlphaBeta(&ab, b, depth, make([]string, 0, 20), game.addChildren)
	fmt.Printf("alphabeta: nodes=%d best score=%v move: %s path: %s\n", ab.nodes, score, move, path)
}

func cmdPlay(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: play fromto\n")
		return
	}
	move := tokens[1]

	// valid move?
	last := len(game.history) - 1
	b := game.history[last]
	if !b.validMove(move) {
		fmt.Printf("not a valid move: %s\n", move)
		return
	}

	for _, c := range b.generateChildren(nil) {
		if c.lastMove == move {
			game.history = append(game.history, c)
			break
		}
	}

	/*
		if len(move) < 4 || len(move) > 5 {
			fmt.Printf("usage: bad move length=%d: '%s'\n", len(move), move)
			return
		}
		from := strings.ToLower(move[:2])
		to := strings.ToLower(move[2:4])

		fmt.Printf("[%s][%s]\n", from, to)

		// from
		if len(from) != 2 {
			fmt.Printf("bad source format: [%s]\n", from)
			return
		}
		if from[0] < 'a' || from[0] > 'h' {
			fmt.Printf("bad source column letter: [%s]\n", from)
		}
		if from[1] < '1' || from[1] > '8' {
			fmt.Printf("bad source row digit: [%s]\n", from)
		}

		// to
		if len(to) != 2 {
			fmt.Printf("bad target format: [%s]\n", to)
			return
		}
		if to[0] < 'a' || to[0] > 'h' {
			fmt.Printf("bad target column letter: [%s]\n", to)
		}
		if to[1] < '1' || to[1] > '8' {
			fmt.Printf("bad target row digit: [%s]\n", to)
		}

		//b := game.history[len(game.history)-1]                        // will update a copy
		p := b.delPiece(location(from[1]-'1'), location(from[0]-'a')) // take piece from board

		if len(move) > 4 {
			// promotion
			promotion := move[4]
			kind := pieceKindFromLetter(rune(promotion))
			if kind != pieceNone {
				p = piece(b.turn<<3) + kind
			}
		}

		b.addPiece(location(to[1]-'1'), location(to[0]-'a'), p) // put piece on board
		b.turn = colorInverse(b.turn)                           // switch color
		b.lastMove = fmt.Sprintf("%s %s", from, to)             // record move

		game.history = append(game.history, b) // append to history
	*/
}

func cmdPerft(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: perft depth\n")
		return
	}
	depth := tokens[1]
	d, errConv := strconv.Atoi(depth)
	if errConv != nil {
		fmt.Printf("bad depth: %s: %v", depth, errConv)
		return
	}

	last := len(game.history) - 1
	b := game.history[last]
	//buf := make([]board, 0, 100)
	buf := []board(nil)
	children := b.generateChildren(buf)

	fmt.Printf("perft depth=%d\n", d)

	total := int64(len(children))
	var nodes int64
	for _, c := range children {
		n, t := perft(c, d, buf)
		fmt.Printf("%s nodes=%d total_nodes=%d\n", c.lastMove, n, t)
		nodes += n
		total += t
	}

	fmt.Printf("perft depth=%d nodes=%d total_nodes=%d\n", d, nodes, total)

	if d < len(testPerftTable) {
		expected := testPerftTable[d+1]
		if expected != nodes {
			fmt.Printf("perft depth=%d nodes=%d expected=%d WRONG\n", d, nodes, expected)
		} else {
			fmt.Printf("perft depth=%d nodes=%d expected=%d ok\n", d, nodes, expected)
		}
	}
}

var testPerftTable = []int64{0, 20, 400, 8902, 197281, 4865609, 119060324, 3195901860}

func perft(b board, depth int, buf []board) (int64, int64) {
	if depth < 1 {
		return 0, 0
	}
	children := b.generateChildren(buf)
	moves := int64(len(children))
	if depth == 1 {
		return moves, moves
	}
	var nodes int64
	for _, c := range children {
		n, total := perft(c, depth-1, buf)
		nodes += n
		moves += total
	}
	return nodes, moves
}

func cmdReset(cmds []command, game *gameState, tokens []string) {
	game.loadFromString(builtinBoard)
}

func cmdSearch(cmds []command, game *gameState, tokens []string) {
	begin := time.Now()

	availTime := 20 * time.Second

	if len(tokens) > 1 {
		a, errParse := time.ParseDuration(tokens[1])
		if errParse != nil {
			fmt.Printf("search: bad duration: '%s': %v\n", tokens[1], errParse)
			return
		}
		availTime = a
	}

	deadline := begin.Add(availTime / 5)

	var bestDepth int
	var bestScore float32
	var bestMove string

LOOP:
	for depth := 1; ; depth++ {
		fmt.Printf("search depth=%d avail=%v remain=%v\n", depth, availTime, time.Until(deadline))
		if deadline.Before(time.Now()) {
			fmt.Printf("search depth=%d: timeout\n", depth)
			break
		}
		last := len(game.history) - 1
		b := game.history[last]
		ab := alphaBetaState{showSearch: false, deadline: deadline}
		score, move, path := rootAlphaBeta(&ab, b, depth, make([]string, 0, 20), game.addChildren)
		if ab.cancelled {
			fmt.Printf("search depth=%d: timeout - cancelled\n", depth)
			break
		}
		fmt.Printf("search depth=%d: nodes=%d best score=%v move: %s path: %s\n", depth, ab.nodes, score, move, path)
		bestDepth = depth
		bestScore = score
		bestMove = move
		switch move {
		case "checkmated", "checkmate", "draw":
			break LOOP
		}
	}

	fmt.Printf("search: best depth=%d score=%v move=%s elapsed=%v\n", bestDepth, bestScore, bestMove, time.Since(begin))
}

func cmdSwitch(cmds []command, game *gameState, tokens []string) {
	b := &game.history[len(game.history)-1] // will update in place
	b.turn = colorInverse(b.turn)           // switch color
}

func cmdUndo(cmds []command, game *gameState, tokens []string) {
	if len(game.history) < 2 {
		return
	}
	game.history = game.history[:len(game.history)-1]
}
