package main

import (
	"fmt"
	"strconv"
	"strings"
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

	total := len(children)
	var nodes int
	for _, c := range children {
		//moves1 := perft(c, d-1)
		n, t := perft(c, d, buf)
		fmt.Printf("%s nodes=%d total_nodes=%d\n", c.lastMove, n, t)
		nodes += n
		total += t
	}

	fmt.Printf("perft depth=%d nodes=%d total_nodes=%d\n", d, nodes, total)
}

func perft(b board, depth int, buf []board) (int, int) {
	if depth < 1 {
		return 0, 0
	}
	children := b.generateChildren(buf)
	moves := len(children)
	if depth == 1 {
		return moves, moves
	}
	var nodes int
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
