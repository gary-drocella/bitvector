package bitvector

import (
	"errors"
)

type Bitvector struct {
	Vector []uint8
}

func (Bv *Bitvector) Init() {
	Bv.Vector = make([]uint8, 1<<7)
}

func (Bv *Bitvector) SetBit(index uint64, value uint8) error {
	if value > 1 {
		return errors.New("value is greater than 1")
	}

	BlockIndex := index / 8
	BlockBitIndex := index % 8

	if BlockIndex+BlockBitIndex > uint64(cap(Bv.Vector)) {
		var OrigSz uint64 = uint64(cap(Bv.Vector))
		var CurrSz = OrigSz

		for BlockIndex+BlockBitIndex >= CurrSz {
			CurrSz = 2 * CurrSz
		}

		VectorSlice := make([]uint8, CurrSz-OrigSz)
		Bv.Vector = append(Bv.Vector, VectorSlice...)
	}

	Block := Bv.Vector[BlockIndex]

	if value == 1 {
		Block |= 1 << BlockBitIndex
	} else {
		Block &= ^(1 << BlockBitIndex)
	}

	Bv.Vector[BlockIndex] = Block

	return nil
}

func (Bv *Bitvector) GetBit(index uint64) (uint8, error) {

	BlockIndex := index / 8
	BlockBitIndex := index % 8

	Block := Bv.Vector[BlockIndex]

	Block &= 1 << BlockBitIndex

	var Value uint8 = 0

	if Block > 0 {
		Value = 1
	}

	return Value, nil
}

func (Bv *Bitvector) Clean() error {
	finished := false
	CurrSize := uint64(cap(Bv.Vector))

	for CurrSize > 128 && !finished {
		for i := uint64(CurrSize / 2); i < CurrSize; i++ {
			if Bv.Vector[i] > 0 {
				finished = true
			}
		}

		if !finished {
			CurrSize /= 2
			Bv.Vector = append([]uint8(nil), Bv.Vector[:CurrSize]...)
		}
	}

	return nil
}
