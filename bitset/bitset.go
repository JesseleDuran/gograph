package bitset

import "math/big"

type Bitset struct {
	*big.Int
}

func NewBigInt() Bitset {
	return Bitset{
		Int: new(big.Int),
	}
}

func (b Bitset) Exists(i int32) bool {
	return b.Int.Bit(int(i)) == 1
}

func (b Bitset) Set(i int32, value bool) {
	if value {
		b.Int.SetBit(b.Int, int(i), 1)
	} else {
		b.Int.SetBit(b.Int, int(i), 0)
	}
}

func (b Bitset) Len() int {
	return b.Int.BitLen()
}
