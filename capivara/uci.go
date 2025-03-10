package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type uciCommand struct {
	name string
	call func(game *gameState, tokens []string)
}

var tableUci = []uciCommand{
	{"uci", uciCmdUci},
	{"isready", uciCmdIsReady},
	{"position", uciCmdPosition},
	{"quit", uciCmdQuit},
	{"go", uciCmdGo},
}

func uciCmdUci(_ *gameState, _ []string) {
	fmt.Println("id name Capivara", fullVersion())
	fmt.Println("id author https://github.com/udhos/capivara")
	fmt.Println("uciok")
}

func uciCmdIsReady(_ *gameState, _ []string) {
	fmt.Println("readyok")
}

func uciCmdQuit(game *gameState, _ []string) {
	game.println("good bye")
	os.Exit(0)
}

func uciCmdPosition(game *gameState, tokens []string) {
	if len(tokens) < 2 {
		return
	}

	var moves []string

	switch tokens[1] {
	case "startpos":
		// position startpos moves e2e4 c7c5

		game.loadFromString(builtinBoard) // reset board

		if len(tokens) < 3 {
			return
		}
		if tokens[2] != "moves" {
			return
		}
		moves = tokens[3:]

	case "fen":
		// position fen r1k4r/p2nb1p1/2b4p/1p1n1p2/2PP4/3Q1NB1/1P3PPP/R5K1 b -    c3 0 19
		// position fen nbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR         w KQkq -  0 1

		fenTokens := tokens[2:]

		var fen []string
		for i, t := range fenTokens {
			if t == "moves" {
				moves = fenTokens[i+1:]
				break
			}
			fen = append(fen, t)
		}

		game.println(fmt.Sprintf("position fen: %v", fen))

		game.loadFromFen(fen)

	default:
		return
	}

	game.println(fmt.Sprintf("position moves: %v", moves))

	// play every move
	for _, m := range moves {
		if errPlay := game.play(m); errPlay != nil {
			game.println(fmt.Sprintf("play error: %v", errPlay))
			return
		}
	}

	game.println(fmt.Sprintf("played: %v", moves))
}

func uciCmdGo(game *gameState, tokens []string) {

	// go wtime 300000 btime 300000 winc 0 binc 0

	game.println(fmt.Sprintf("version %s", shortVersion()))

	game.println(fmt.Sprintf("go: %v", tokens))

	avail := 30 * time.Second // just a default

	turn := game.history[len(game.history)-1].turn

	game.println(fmt.Sprintf("turn: %s", turn.name()))

	var timeLabel string
	if turn == colorWhite {
		timeLabel = "wtime"
	} else {
		timeLabel = "btime"
	}

	var perMove bool

	for i := 1; i < len(tokens)-1; i++ {
		t := tokens[i]

		if t == "movetime" {
			// found per-move time
			avail = parseTime(game, tokens[i+1])
			perMove = true
			break
		}

		if t == timeLabel {
			// found remaining time
			avail = parseTime(game, tokens[i+1])
			break
		}
	}

	game.println(fmt.Sprintf("available time: %v", avail))

	var bestMove string

	if perMove {
		bestMove = game.searchPerMove(avail, avail)
	} else {
		bestMove = game.search(avail)
	}

	fmt.Println("bestmove", bestMove)
}

func parseTime(game *gameState, t string) time.Duration {
	v, errConv := strconv.Atoi(t)
	if errConv != nil {
		game.println(fmt.Sprintf("error: %s: %v", t, errConv))
	}
	avail := time.Duration(v) * time.Millisecond
	return avail
}
