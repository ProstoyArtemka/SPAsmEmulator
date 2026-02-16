package data

import (
	"math/bits"
)

const (
	R0 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	SP
	LR
	PC
	PSR
)

const (
	N = iota
	Z
	C
	V
	Q
)

type Register struct {
	Value int32
}

var NIL_REGISTER = Register{}

func (r *Register) SetBit(index int, val int) {

	if val == 0 {
		r.Value &^= (1 << index)

		return
	}

	r.Value |= (1 << index)

}

func (r *Register) GetBit(index int) byte {

	uRegister := uint32(r.Get())

	if (uRegister>>index)&1 != 0 {
		return 1
	}

	return 0
}

func (r *Register) Set(val int32) {
	r.Value = val
}

func (r *Register) Add(val int32) (int32, int32) {

	sum, carry := bits.Add32(uint32(r.Value), uint32(val), 0)

	r.Value = int32(sum)

	return int32(sum), int32(carry)
}

func (r *Register) Sub(val int32) (int32, int32) {

	sum, carry := bits.Sub32(uint32(r.Value), uint32(val), 0)

	r.Value = int32(sum)

	return int32(sum), int32(carry)
}

func (r *Register) Get() int32 {
	return r.Value
}

func (r *Register) Increment() (int32, int32) {

	return r.Add(1)

}
