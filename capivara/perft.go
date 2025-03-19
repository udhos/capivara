package main

var testPerftTable = []int64{0, 20, 400, 8902, 197281, 4865609, 119060324, 3195901860}

func perft(b board, depth int, buf *boardPool) (int64, int64) {
	if depth < 1 {
		return 0, 0
	}
	const pruneRepetition = false
	countChildren, _ := b.generateChildren(buf, pruneRepetition)
	//log.Printf("perft +++ depth=%d children=%d pool=%d", depth, countChildren, len(buf.pool))
	moves := int64(countChildren)
	if depth == 1 {
		buf.drop(countChildren)
		//log.Printf("perft ---- depth=%d children=%d pool=%d", depth, countChildren, len(buf.pool))
		return moves, moves
	}
	var nodes int64
	lastChildren := buf.pool[len(buf.pool)-countChildren:]
	for _, c := range lastChildren {
		n, total := perft(c, depth-1, buf)
		nodes += n
		moves += total
	}
	buf.drop(countChildren)
	//log.Printf("perft ---- depth=%d children=%d pool=%d", depth, countChildren, len(buf.pool))
	return nodes, moves
}
