package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func (game *gameState) bookLookup() string {

	history := game.history
	if history[0].lastMove.isNull() {
		history = history[1:]
	}

	var moves []string
	for _, m := range history {
		moves = append(moves, m.lastMove.String())
	}
	position := strings.Join(moves, " ")
	game.println(fmt.Sprintf("bookLookup: position: [%s]", position))

	if moves, found := book[position]; found {
		return game.bookPick(position, moves)
	}

	return "" // not found
}

func (game *gameState) bookPick(position string, moves []bookMove) string {

	if len(moves) < 1 {
		game.println(fmt.Sprintf("bookLookup: position missing moves: [%s]", position))
		return "" // not found
	}

	if len(moves) == 1 {
		return moves[0].move // single move
	}

	// multiple moves

	var sum int
	for _, m := range moves {
		w := m.getWeight(game, position)
		sum += w
	}

	r := int(rand.Int31n(int32(sum))) // 0..sum-1

	// 2 3 4
	// 2 5 9
	// sum = 9
	// 0..8
	// 0..1 -> 2
	// 2..4 -> 3
	// 5..8 -> 4

	var rs int // running sum
	for i, m := range moves {
		w := m.getWeight(game, position)
		rs += w
		if rs > r {
			game.println(fmt.Sprintf("bookPick: position=[%s] rand=%d runSum=%d index=%d move=%s", position, r, rs, i, m.move))
			return m.move // found
		}
	}

	game.println(fmt.Sprintf("bookPick: position=[%s] not reached - ugh", position))

	return moves[0].move
}

type stringReader interface {
	ReadString(delim byte) (string, error) // Example: bufio.Reader
}

func loadBook(reader stringReader) {

	var lineCount int

LOOP:
	for {
		lineCount++
		line, errRead := reader.ReadString('\n')
		switch errRead {
		case io.EOF:
			loadLine(lineCount, line) // last line
			break LOOP
		case nil:
			loadLine(lineCount, line)
		default:
			loadLine(lineCount, line)
			log.Printf("loadBook: read error at line=%d: %v", lineCount, errRead)
		}
	}

	log.Printf("loadBook: lines=%d bookSize=%d", lineCount, len(book))
}

func loadLine(count int, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	entry := strings.SplitN(line, ":", 2)
	if len(entry) != 2 {
		log.Printf("loadLine: missing position:move at line=%d: %s", count, line)
		return
	}
	position := strings.TrimSpace(entry[0])

	moves := strings.Split(strings.TrimSpace(entry[1]), ",")

	for _, moveWeight := range moves {

		w := 1 // default weight

		mw := strings.SplitN(strings.TrimSpace(moveWeight), " ", 2)

		moveStr := strings.TrimSpace(mw[0])

		if len(mw) > 1 {
			value, errConv := strconv.Atoi(strings.TrimSpace(mw[1]))
			if errConv != nil {
				log.Printf("loadLine: bad move weight at line=%d: %s: %v", count, line, errConv)
			} else {
				w = value
			}
		}

		log.Printf("loadLine: line=%d: position=[%s] move=%s weight=%d", count, position, moveStr, w)

		book[position] = append(book[position], bookMove{move: moveStr, weight: w})

	}
}

const defaultBook = `
: e2e4 2, d2d4, b1f3
e2e4: c7c5 2, e7e5, e7e6
e2e4 c7c5: g1f3 2, b1c3, c2c3
`

var book = map[string][]bookMove{
	//"":     {{move: "e2e4", weight: 1}, {move: "d2d4", weight: 1}, {move: "b1f3", weight: 1}},
	//"e2e4": {{move: "c7c5", weight: 1}},
}

type bookMove struct {
	move   string
	weight int
}

func (m bookMove) getWeight(g *gameState, position string) int {
	w := m.weight
	if w < 1 {
		g.println(fmt.Sprintf("bookLookup: bad weight=%d for move=%s position: [%s]", w, m.move, position))
		w = 1
	}
	return w
}
