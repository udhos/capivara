package main

type boardPool struct {
	pool []board
}

var defaultBoardPool = newPool()

func newPool() *boardPool {
	return newPoolSize(1000)
}

func newPoolSize(size int) *boardPool {
	return &boardPool{pool: make([]board, 0, size)}
}

func (bp *boardPool) push(b *board) {
	bp.pool = append(bp.pool, *b)
}

func (bp *boardPool) drop(n int) {
	bp.pool = bp.pool[:len(bp.pool)-n]
}

func (bp *boardPool) reset() {
	bp.pool = bp.pool[:0]
}
