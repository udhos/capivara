package main

import (
	"fmt"
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
	{"go", uciCmdGo},
}

func uciCmdUci(game *gameState, tokens []string) {
	fmt.Println("id name Capivara", fullVersion())
	fmt.Println("id author https://github.com/udhos/capivara")
	fmt.Println("uciok")
}

func uciCmdIsReady(game *gameState, tokens []string) {
	fmt.Println("readyok")
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

	game.println(fmt.Sprintf("go: %v", tokens))

	avail := 30 * time.Second // just a default

	turn := game.history[len(game.history)-1].turn

	var timeLabel string
	if turn == colorWhite {
		timeLabel = "wtime"
	} else {
		timeLabel = "btime"
	}

	for i := 1; i < len(tokens)-1; i++ {
		t := tokens[i]
		if t == timeLabel {
			tt := tokens[i+1]
			v, errConv := strconv.Atoi(tt)
			if errConv != nil {
				game.println(fmt.Sprintf("error: %s %s: %v", timeLabel, tt, errConv))
			}
			avail = time.Duration(v) * time.Millisecond
			break
		}
	}

	game.println(fmt.Sprintf("available time: %v", avail))

	bestMove := game.search(avail)

	fmt.Println("bestmove", bestMove)
}
