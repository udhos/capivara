package main

import "testing"

func TestColorToSignal(t *testing.T) {

	sigWhite := colorToSignal(colorWhite)
	if sigWhite != 1 {
		t.Errorf("colorToSignal(white) != 1 (got %d)", sigWhite)
	}

	sigBlack := colorToSignal(colorBlack)
	if sigBlack != -1 {
		t.Errorf("colorToSignal(black) != -1 (got %d)", sigBlack)
	}
}

func TestColorInverse(t *testing.T) {

	whiteInverse := colorInverse(colorWhite)
	if whiteInverse != colorBlack {
		t.Errorf("colorInverse(white) != black (got %d)", whiteInverse)
	}

	blackInverse := colorInverse(colorBlack)
	if blackInverse != colorWhite {
		t.Errorf("colorInverse(black) != white (got %d)", blackInverse)
	}
}
