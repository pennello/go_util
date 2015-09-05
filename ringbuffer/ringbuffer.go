// chris 00515

// Package ringbuffer implements an append-only ring buffer.
//
// This is also known as a circular or cyclic buffer.  It uses a single,
// fixed-size buffer as if it were connected end-to-end.
package ringbuffer

// B is the type of the append-only ring buffer.  This data structure is
// also known as a circular buffer or cyclic buffer.  It uses a single,
// fixed-size buffer as if it were connected end-to-end.
type B struct {
	// The size of the ring buffer (the value passed into New).
	Size int

	// The underlying buffer.
	buf []byte

	// The write position in the buffer.
	pos int

	// Indicates whether we've written enough to the ring buffer to
	// wrap around again.  True iff at least Size bytes have been
	// written.
	wrapped bool
}

// New allocates a new ring buffer of the specified size and returns a
// pointer to it.  If size is non-positive, New panics.
func New(size int) *B {
	if size < 1 {
		panic("non-positive size")
	}

	return &B{
		Size: size,
		buf:  make([]byte, size),
		pos:  0,
	}
}

// remaining returns the remaining space in b.buf until we have to wrap.
func (b *B) remaining() int {
	return b.Size - b.pos
}

// Write writes len(p) bytes from p into the ring buffer.  It returns
// the number of bytes written from p, which is always len(p); and err,
// which is always nil.  If len(p) > b.Size, then only the final len(p)
// - b.Size bytes will remain in the ring buffer.
func (b *B) Write(p []byte) (n int, err error) {
	n = len(p)

	if len(p) > b.Size {
		p = p[len(p)-b.Size:]
	}

	// Now len(p) <= b.Size.

	r := b.remaining()

	if len(p) <= r {
		copy(b.buf[b.pos:], p)
	} else {
		copy(b.buf[b.pos:], p[:r])
		copy(b.buf, p[r:])
	}

	b.pos += len(p)
	if b.pos >= b.Size {
		b.wrapped = true
		b.pos -= b.Size
	}

	return n, nil
}

// Bytes returns the contents of the ring buffer in logical order (i.e.,
// not as it's actually stored in the underlying buffer).
//
// If at least b.Size bytes have been written into the ring buffer, then
// Bytes will allocate a new byte slice of length b.Size, and copy into
// it the contents of the ring buffer.
//
// Otherwise, if less than b.Size bytes have been written into the ring
// buffer, then a slice containing only what's been written (potentially
// of size less than b.Size) will be returned.  The slice will share
// storage with the underlying buffer used by the ring buffer.
func (b *B) Bytes() []byte {
	if !b.wrapped {
		return b.buf[:b.pos]
	}
	// We _could_ try to re-use b.buf here, slicing it and
	// appending, but we'd have to allocate a new array anyway!
	bytes := make([]byte, b.Size)
	copy(bytes, b.buf[b.pos:])
	copy(bytes[b.remaining():], b.buf[:b.pos])
	return bytes
}
