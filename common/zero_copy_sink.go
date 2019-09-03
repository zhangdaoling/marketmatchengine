package common

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ZeroCopyBytes struct {
	buf []byte
}

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (z *ZeroCopyBytes) tryGrowByReslice(n int) (int, bool) {
	if l := len(z.buf); n <= cap(z.buf)-l {
		z.buf = z.buf[:l+n]
		return l, true
	}
	return 0, false
}

const maxInt = int(^uint(0) >> 1)

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (z *ZeroCopyBytes) grow(n int) int {
	// Try to grow by means of a reslice.
	if i, ok := z.tryGrowByReslice(n); ok {
		return i
	}

	l := len(z.buf)
	c := cap(z.buf)
	if c > maxInt-c-n {
		panic(ErrTooLarge)
	}
	// Not enough space anywhere, we need to allocate.
	buf := makeSlice(2*c + n)
	copy(buf, z.buf)
	z.buf = buf[:l+n]
	return l
}

func (z *ZeroCopyBytes) WriteBytes(p []byte) {
	data := z.NextBytes(uint64(len(p)))
	copy(data, p)
}

func (z *ZeroCopyBytes) Size() uint64 { return uint64(len(z.buf)) }

func (z *ZeroCopyBytes) NextBytes(n uint64) (data []byte) {
	m, ok := z.tryGrowByReslice(int(n))
	if !ok {
		m = z.grow(int(n))
	}
	data = z.buf[m:]
	return
}

// Backs up a number of bytes, so that the next call to NextXXX() returns data again
// that was already returned by the last call to NextXXX().
func (z *ZeroCopyBytes) BackUp(n uint64) {
	l := len(z.buf) - int(n)
	z.buf = z.buf[:l]
}

func (z *ZeroCopyBytes) WriteUint8(data uint8) {
	buf := z.NextBytes(1)
	buf[0] = data
}

func (z *ZeroCopyBytes) WriteByte(c byte) {
	z.WriteUint8(c)
}

func (z *ZeroCopyBytes) WriteBool(data bool) {
	if data {
		z.WriteByte(1)
	} else {
		z.WriteByte(0)
	}
}

func (z *ZeroCopyBytes) WriteUint16(data uint16) {
	buf := z.NextBytes(2)
	binary.LittleEndian.PutUint16(buf, data)
}

func (z *ZeroCopyBytes) WriteUint32(data uint32) {
	buf := z.NextBytes(4)
	binary.LittleEndian.PutUint32(buf, data)
}

func (z *ZeroCopyBytes) WriteUint64(data uint64) {
	buf := z.NextBytes(8)
	binary.LittleEndian.PutUint64(buf, data)
}

func (z *ZeroCopyBytes) WriteInt64(data int64) {
	z.WriteUint64(uint64(data))
}

func (z *ZeroCopyBytes) WriteInt32(data int32) {
	z.WriteUint32(uint32(data))
}

func (z *ZeroCopyBytes) WriteInt16(data int16) {
	z.WriteUint16(uint16(data))
}

func (z *ZeroCopyBytes) WriteVarBytes(data []byte) (size uint64) {
	l := uint64(len(data))
	size = z.WriteVarUint(l) + l

	z.WriteBytes(data)
	return
}

func (z *ZeroCopyBytes) WriteString(data string) (size uint64) {
	return z.WriteVarBytes([]byte(data))
}

func (z *ZeroCopyBytes) WriteVarUint(data uint64) (size uint64) {
	buf := z.NextBytes(9)
	if data < 0xFD {
		buf[0] = uint8(data)
		size = 1
	} else if data <= 0xFFFF {
		buf[0] = 0xFD
		binary.LittleEndian.PutUint16(buf[1:], uint16(data))
		size = 3
	} else if data <= 0xFFFFFFFF {
		buf[0] = 0xFE
		binary.LittleEndian.PutUint32(buf[1:], uint32(data))
		size = 5
	} else {
		buf[0] = 0xFF
		binary.LittleEndian.PutUint64(buf[1:], uint64(data))
		size = 9
	}

	z.BackUp(9 - size)
	return
}

// NewReader returns a new ZeroCopyBytes reading from b.
func NewZeroCopyBytes(b []byte) *ZeroCopyBytes {
	if b == nil {
		b = make([]byte, 0, 512)
	}
	return &ZeroCopyBytes{b}
}

func (z *ZeroCopyBytes) Bytes() []byte { return z.buf }

func (z *ZeroCopyBytes) Reset() { z.buf = z.buf[:0] }

var ErrTooLarge = errors.New("bytes.Buffer: too large")

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(bytes.ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
