#!/usr/bin/env python3

import fileinput
import chess.pgn

game = chess.pgn.Game()
node = game

print("reading UCI moves from stdin...")

for line in fileinput.input():
    moves = line.split()
    for m in moves:
        mv = m.strip(" \r\n\t")
        print("adding", mv)
        node = node.add_variation(chess.Move.from_uci(mv))

print("reading UCI moves from stdin...done")

print(game)