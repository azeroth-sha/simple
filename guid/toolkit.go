package guid

import (
	"math/big"
	"sync"
)

var guidPool = &sync.Pool{New: func() any { return new(big.Int) }}

/*
  Package Private functions
*/

func getInt() *big.Int {
	return guidPool.Get().(*big.Int)
}

func putInt(bInt *big.Int) {
	bInt.SetInt64(0)
	guidPool.Put(bInt)
}
