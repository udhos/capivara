package main

import (
	"testing"
)

type expect bool

const (
	expectSuccess expect = true
	expectError   expect = false
)

type positionTest struct {
	name         string
	position     string
	expectResult expect
}

var testPositionTable = []positionTest{
	{"e2e4", "e2e4", expectSuccess},
	{"e2e5", "e2e5", expectError},
	{"e2e4 e7e5", "e2e4 e7e5", expectSuccess},
	{"e2e4 c7c5", "e2e4 c7c5", expectSuccess},
	{"e2e4 e2e4", "e2e4 e2e4", expectError},
	{"e2e4 e8e7", "e2e4 e8e7", expectError},
}

func TestPosition(t *testing.T) {

	for _, data := range testPositionTable {

		tmpG := newGame()
		tmp := &tmpG
		tmp.loadFromString(builtinBoard)
		var errTmp error
		_, errTmp = tmp.validatePosition(data.position)

		if result := expect(errTmp == nil); result != data.expectResult {
			t.Errorf("%s: position=[%s] expected=%v got=%v\n",
				data.name, data.position, data.expectResult, result)
		}
	}
}
