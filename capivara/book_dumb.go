package main

const defaultBook = `
# format for specific position
# this format adds one or more moves as responses to a specific position.
#
# position: move [weigth] [... , move [weight]]
#
# examples:
#
# e2e4: c7c5         # adds c7c5 as response for e2e4
# e2e4: c7c5 2, e7e5 # adds two responses: c7c5 with weight 2 and e7e5 with weight 1

: e2e4 2, d2d4, g1f3
e2e4: c7c5 2, e7e5, e7e6
e2e4 c7c5: g1f3 2, b1c3, c2c3
e2e4 c7c5 g1f3: d7d6, b8c6, e7e6

: e2e4 # dup move

# format for full game
# this format adds all moves from a sequence of moves.
# problem is, all moves are added as strong responses.
#
# position
#
# example:
#
# d2d4 g8f6 c2c4 e7e6

d2d4 g8f6 c2c4 e7e6
d2d4 g8f6 c2c4 e7e6 : b1c3
d2d4 d7d5 c2c4
d2d4 d7d5 c2c4 e7e6 : b1c3
g1f3 d7d5 g2g3 e7e6 f1g2
g1f3 d7d5 g2g3 g8f6 f1g2
e2e4 e7e5 f1c4 g8f6 d2d3
e2e4 e7e5 f1c4 g8f6 d2d3 f8b4 : c2c3
`
