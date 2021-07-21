package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func (game *gameState) bookLookup() string {

	position := game.position()
	game.println(fmt.Sprintf("bookLookup: position: [%s]", position))

	if moves, found := book[position]; found {
		return game.bookPick(position, moves)
	}

	return "" // not found
}

func (game *gameState) position() string {
	history := game.history
	if history[0].lastMove.isNull() {
		history = history[1:]
	}

	var moves []string
	for _, m := range history {
		moves = append(moves, m.lastMove.String())
	}
	position := strings.Join(moves, " ")

	return position
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

func loadBookFromFile(filename string) {
	fmt.Printf("loadBookFromFile: %s\n", filename)
	input, errOpen := os.Open(filename)
	if errOpen != nil {
		fmt.Printf("loadBookFromFile: %s: %v\n", filename, errOpen)
		return
	}
	defer input.Close()
	loadBook(bufio.NewReader(input))
}

var book = map[string][]bookMove{}

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

func loadBook(reader stringReader) {

	book = map[string][]bookMove{}

	var lineCount int

LOOP:
	for {
		lineCount++
		line, errRead := reader.ReadString('\n')
		switch errRead {
		case io.EOF:
			// last line
			loadLine(lineCount, line)
			break LOOP
		case nil:
			if loadLine(lineCount, line) {
				break LOOP
			}
		default:
			fatal := loadLine(lineCount, line)
			log.Printf("loadBook: read error at line=%d: %v", lineCount, errRead)
			if fatal {
				break LOOP
			}
		}
	}

	log.Printf("loadBook: lines=%d bookSize=%d", lineCount, len(book))
}

const errFatal = true
const errNonFatal = false

func loadLine(count int, line string) bool {

	comment := strings.SplitN(line, "#", 2)
	uncomment := strings.TrimSpace(comment[0])

	if uncomment == "" {
		return errNonFatal
	}

	entry := strings.SplitN(uncomment, ":", 2)
	if len(entry) < 1 {
		log.Printf("loadLine: missing position at line=%d: %s", count, line)
		return errNonFatal
	}
	positionMoves := strings.Fields(entry[0])
	position := strings.Join(positionMoves, " ")

	tmpG := newGame()
	tmp := &tmpG
	tmp.loadFromString(builtinBoard)
	var errTmp error
	tmp, errTmp = tmp.validatePosition(position)
	if errTmp != nil {
		log.Printf("loadLine: line=%d: invalid position=[%s]: %v", count, position, errTmp)
		return errFatal
	}

	if len(entry) == 1 {
		return loadGame(count, positionMoves)
	}

	moves := strings.Split(strings.TrimSpace(entry[1]), ",")

	for _, moveWeight := range moves {

		w := 1 // default weight

		mw := strings.SplitN(strings.TrimSpace(moveWeight), " ", 2)

		moveStr := strings.TrimSpace(mw[0])

		tmp, errTmp = tmp.validatePosition(moveStr)
		if errTmp != nil {
			log.Printf("loadLine: line=%d: invalid move position=[%s]: move=%s %v", count, position, moveStr, errTmp)
			return errFatal
		}
		tmp.undo()

		if len(mw) > 1 {
			value, errConv := strconv.Atoi(strings.TrimSpace(mw[1]))
			if errConv != nil {
				log.Printf("loadLine: bad move weight at line=%d: %s: %v", count, line, errConv)
			} else {
				w = value
			}
		}

		//log.Printf("loadLine: line=%d: position=[%s] move=%s weight=%d", count, position, moveStr, w)
		//book[position] = append(book[position], bookMove{move: moveStr, weight: w})
		loadPosition(position, bookMove{move: moveStr, weight: w}, count)
	}

	return errNonFatal
}

func loadPosition(position string, m bookMove, count int) {
	//log.Printf("loadPosition: line=%d: position=[%s] move=%s weight=%d FIXME PREVENT DUP MOVE", count, position, m.move, m.weight)
	book[position] = append(book[position], m)
}

func loadGame(count int, positionMoves []string) bool {
	//full := strings.Join(positionMoves, " ")
	var moves []string
	for _, m := range positionMoves {
		position := strings.Join(moves, " ")
		//log.Printf("loadGame: line=%d: position=[%s] p=[%s] move=%s", count, full, position, m)
		//book[position] = append(book[position], bookMove{move: m, weight: 1})
		loadPosition(position, bookMove{move: m, weight: 1}, count)
		moves = append(moves, m)
	}
	return errNonFatal
}

func (game *gameState) validatePosition(position string) (*gameState, error) {
	moves := strings.Fields(position)
	for _, t := range moves {
		if errPlay := game.play(t); errPlay != nil {
			return game, errPlay // error
		}
	}
	return game, nil // ok
}
