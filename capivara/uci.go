package main

import (
	"fmt"
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
	fmt.Println("id name Capivara")
	fmt.Println("id author https://github.com/udhos/capivara")
	fmt.Println("uciok")
}

func uciCmdIsReady(game *gameState, tokens []string) {
	fmt.Println("readyok")
}

func uciCmdPosition(game *gameState, tokens []string) {
	if tokens[1] == "startpos" && tokens[2] == "moves" {
		m := tokens[3]

		game.loadFromString(builtinBoard)

		if errPlay := game.play(m); errPlay != nil {
			game.println(fmt.Sprintf("play error: %v", errPlay))
			return
		}

		game.println(fmt.Sprintf("played %s", m))
	}
}

func uciCmdGo(game *gameState, tokens []string) {
	game.println(fmt.Sprintf("go: %v", tokens))
}
