# capivara

# Features

* [UCI protocol](http://wbec-ridderkerk.nl/html/UCIProtocol.html)
* [Alpha-beta search](https://www.chessprogramming.org/Alpha-Beta)
* [Iterative deepening](https://www.chessprogramming.org/Iterative_Deepening)
* [Piece-square tables](https://www.chessprogramming.org/Piece-Square_Tables)

# How to build

    git clone https://github.com/udhos/capivara
    cd capivara
    go install ./capivara

# How to play

## Run

Run `capivara`.

    capivara

Then type commands after the prompt `enter command:`.
Commands may be abbreviated. For instance, use `h` for `help`.

## Search best move

Use command `search` to ask for the best move.

    enter command:s

Engine response will look like this:

    search: best depth=4 score=-1 move=d7d6 elapsed=4.000061404s

## Play a move

Use command `play <move>` to make a move.

    enter command:p e2e4

## Help on commands

Use command `help` to get help. 

    enter command:h

# Other Golang Chess Engines

https://github.com/tonyOreglia/glee

https://github.com/zserge/carnatus

https://github.com/ChizhovVadim/CounterGo

https://github.com/WeaselFOSS/Chess

https://github.com/amanjpro/zahak

https://github.com/frankkopp/FrankyGo
