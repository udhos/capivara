package main

import (
	"fmt"
	"log"
)

var nullMove = move{}

type move struct {
	src       location
	dst       location
	promotion piece
}

func newMove(s string) move {
	src := s[:2]
	dst := s[2:4]

	var p piece

	if len(s) > 4 {
		// promotion
		promotion := s[4]
		p = pieceKindFromLetter(rune(promotion))
	}

	m := move{
		src:       8*location(src[1]-'1') + location(src[0]-'a'),
		dst:       8*location(dst[1]-'1') + location(dst[0]-'a'),
		promotion: p,
	}

	log.Printf("NewMove: %s => %s", s, m)

	return m
}

func (m move) isNull() bool {
	return m.src == 0 && m.dst == 0
}

func (m move) rankDelta() int {
	return int(abs(int64(m.dst-m.src) >> 3))
}

func (m move) String() string {
	srcRow := m.src / 8
	srcCol := m.src % 8
	dstRow := m.dst / 8
	dstCol := m.dst % 8
	s := fmt.Sprintf("%c%c%c%c", srcCol+'a', srcRow+'1', dstCol+'a', dstRow+'1')
	if m.promotion != pieceNone {
		s += m.promotion.kindLetterLow()
	}
	return s
}
