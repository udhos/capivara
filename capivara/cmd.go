package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
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
	{"castling", cmdCastling, "castling"},
	{"clear", cmdClear, "erase board"},
	{"dumbbook", cmdDumbBook, "dumb book"},
	{"fen", cmdFen, "load board from FEN"},
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
	{"uci", cmdUci, "start UCI mode"},
}

func cmdClear(cmds []command, game *gameState, tokens []string) {
	*game = newGame()
}

func cmdCastling(cmds []command, game *gameState, tokens []string) {
	last := len(game.history) - 1
	b := &game.history[last]
	b.flags[colorWhite] |= lostCastlingLeft | lostCastlingRight // disable castling for white
	b.flags[colorBlack] |= lostCastlingLeft | lostCastlingRight // disable castling for black
	fmt.Println("castling disabled")
}

func cmdDumbBook(cmds []command, game *gameState, tokens []string) {
	game.dumbBook = !game.dumbBook
	fmt.Println("dumb book:", game.dumbBook)
}

func cmdFen(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: fen FEN-string\n")
		return
	}
	game.loadFromFen(tokens[1:])
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

	children := defaultBoardPool
	children.reset()
	nega := negamaxState{children: children, showSearch: true}

	begin := time.Now()

	score, move, comment := rootNegamax(&nega, b, depth, game.addChildren)

	speed := getSpeed(nega.nodes, begin)

	fmt.Printf("negamax: nodes=%d speed=%v knodes/s best score=%v move=%s (%s)\n", nega.nodes, speed, score, move, comment)
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

	children := defaultBoardPool
	children.reset()
	ab := alphaBetaState{showSearch: true, children: children}

	begin := time.Now()

	score, move, comment := rootAlphaBeta(&ab, b, depth, game.addChildren)

	speed := getSpeed(ab.nodes, begin)

	fmt.Printf("alphabeta: nodes=%d speed=%v knodes/s best score=%v move=%s (%s)\n", ab.nodes, speed, score, move, comment)
}

func getSpeed(nodes int64, begin time.Time) int {
	elap := time.Since(begin)
	return getSpeedElapsed(nodes, elap)
}

func getSpeedElapsed(nodes int64, elap time.Duration) int {
	return int(float64(nodes/1000) / elap.Seconds()) // knodes / s
}

func cmdPlay(cmds []command, game *gameState, tokens []string) {
	if len(tokens) < 2 {
		fmt.Printf("usage: play fromto\n")
		return
	}

	isSep := func(c rune) bool {
		return c == '(' || c == ')'
	}

	for _, t := range tokens[1:] {
		fields := strings.FieldsFunc(t, isSep)
		for _, m := range fields {
			if errPlay := game.play(m); errPlay != nil {
				fmt.Printf("play error: %v\n", errPlay)
				return
			}
		}
	}
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

	perftBegin := time.Now()

	last := len(game.history) - 1
	b := game.history[last]
	//buf := make([]board, 0, 100)
	//buf := []board(nil)
	children := defaultBoardPool
	children.reset()
	countChildren := b.generateChildren(children)

	fmt.Printf("perft depth=%d\n", d)

	total := int64(countChildren)
	var nodes int64
	for _, c := range children.pool {
		begin := time.Now()
		n, t := perft(c, d, children)
		elap := time.Since(begin)
		speed := getSpeedElapsed(t, elap)
		fmt.Printf("%s nodes=%d total_nodes=%d elapsed=%v speed=%v knodes/s\n", c.lastMove, n, t, elap, speed)
		nodes += n
		total += t
	}

	perftElap := time.Since(perftBegin)
	perftSpeed := getSpeedElapsed(total, perftElap)

	fmt.Printf("perft depth=%d nodes=%d total_nodes=%d elapsed=%v speed=%v knodes/s\n", d, nodes, total, perftElap, perftSpeed)

	if d < len(testPerftTable) {
		expected := testPerftTable[d+1]
		if expected != nodes {
			fmt.Printf("perft depth=%d nodes=%d expected=%d WRONG\n", d, nodes, expected)
		} else {
			fmt.Printf("perft depth=%d nodes=%d expected=%d ok\n", d, nodes, expected)
		}
	}
}

func cmdReset(cmds []command, game *gameState, tokens []string) {
	game.loadFromString(builtinBoard)
}

func cmdSearch(cmds []command, game *gameState, tokens []string) {
	availTime := 60 * time.Second

	if len(tokens) > 1 {
		a, errParse := time.ParseDuration(tokens[1])
		if errParse != nil {
			fmt.Printf("search: bad duration: '%s': %v\n", tokens[1], errParse)
			return
		}
		availTime = a
	}

	game.search(availTime)
}

func (game *gameState) search(availTime time.Duration) string {

	if game.dumbBook {
		best := game.bookLookup()
		if best != "" {
			game.println(fmt.Sprintf("dumb book best move: %s", best))
			return best // found
		}
	}

	begin := time.Now()

	deadline := begin.Add(availTime / 20)

	var bestDepth int
	var bestScore float32
	var bestMove move
	var bestComment string

	if game.cpuprofile != "" {
		f, err := os.Create(game.cpuprofile)
		if err != nil {
			log.Printf("cpuprofile: %v", err)
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

LOOP:
	for depth := 1; ; depth++ {
		game.print(fmt.Sprintf("search depth=%d avail=%v remain=%v\n", depth, availTime, time.Until(deadline)))
		depthBegin := time.Now()
		if deadline.Before(depthBegin) {
			game.print(fmt.Sprintf("search depth=%d: timeout\n", depth))
			break
		}

		children := defaultBoardPool
		children.reset()
		ab := alphaBetaState{showSearch: false, deadline: deadline, children: children}

		last := len(game.history) - 1
		b := game.history[last]
		score, move, comment := rootAlphaBeta(&ab, b, depth, game.addChildren)
		if ab.cancelled {
			game.print(fmt.Sprintf("search depth=%d: timeout - cancelled\n", depth))
			break
		}

		speed := getSpeed(ab.nodes, depthBegin)

		game.print(fmt.Sprintf("search depth=%d: nodes=%d speed=%v knodes/s best score=%v move=%s (%s)\n", depth, ab.nodes, speed, score, move, comment))
		bestDepth = depth
		bestScore = score
		bestMove = move
		bestComment = comment
		if ab.singleChildren {
			game.print(fmt.Sprintf("search depth=%d: move=%s single move\n", depth, move))
			break
		}
		switch comment {
		case "checkmated", "checkmate", "draw":
			break LOOP
		}
		if bestScore == alphabetaMax {
			game.print(fmt.Sprintf("search depth=%d: nodes=%d best score=%v move: %s found checkmate\n", depth, ab.nodes, score, move))
			break
		}
	}

	game.println(fmt.Sprintf("search: best depth=%d score=%v move=%s elapsed=%v", bestDepth, bestScore, bestMove, time.Since(begin)))

	if bestMove.isNull() {
		return bestComment
	}

	return bestMove.String()
}

func cmdSwitch(cmds []command, game *gameState, tokens []string) {
	b := &game.history[len(game.history)-1] // will update in place
	b.turn = colorInverse(b.turn)           // switch color
}

func cmdUci(cmds []command, game *gameState, tokens []string) {
	uciCmdUci(game, tokens)
	game.uci = true
}

func cmdUndo(cmds []command, game *gameState, tokens []string) {
	if len(game.history) < 2 {
		return
	}
	game.history = game.history[:len(game.history)-1]
}
