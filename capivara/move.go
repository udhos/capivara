package main

import (
	"fmt"
)

var nullMove = move{}

type move struct {
	src       location
	dst       location
	promotion piece
}

func newMove(s string) (move, error) {
	if len(s) < 4 {
		return nullMove, fmt.Errorf("newMove: bad move length(%s)=%d < 4", s, len(s))
	}

	src := s[:2]
	dst := s[2:4]

	var p piece

	if len(s) > 4 {
		// promotion
		promotion := s[4]
		p = pieceKindFromLetter(rune(promotion))
		if p == pieceNone {
			return nullMove, fmt.Errorf("newMove: bad move promotion: %s", s)
		}
	}

	if src[0] < 'a' || src[0] > 'h' {
		return nullMove, fmt.Errorf("newMove: bad move source column: %s", s)
	}

	if src[1] < '1' || src[1] > '8' {
		return nullMove, fmt.Errorf("newMove: bad move source rank: %s", s)
	}

	if dst[0] < 'a' || dst[0] > 'h' {
		return nullMove, fmt.Errorf("newMove: bad move destination column: %s", s)
	}

	if dst[1] < '1' || dst[1] > '8' {
		return nullMove, fmt.Errorf("newMove: bad move destination rank: %s", s)
	}

	m := move{
		src:       8*location(src[1]-'1') + location(src[0]-'a'),
		dst:       8*location(dst[1]-'1') + location(dst[0]-'a'),
		promotion: p,
	}

	return m, nil
}

func (m move) isNull() bool {
	return m.src == 0 && m.dst == 0
}

func (m move) rankDelta() int {
	return int(abs(int64(m.dst/8 - m.src/8)))
}

func (m move) String() string {
	if m.isNull() {
		return ""
	}
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
