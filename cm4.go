package boundedkv

import (
	"hash/maphash"
)

// CM4 is a small conservative-update count-min sketch implementation with 4-bit counters.
// This provides an estimation of the relative frequency of items.
// The zero value is usable, and tuned for 32-bits numbers.
type CM4 struct {
	nyb [width * depth / 2]byte
	h   maphash.Hash
	ins uint16
}

const width = 32
const depth = 4

func (c *CM4) Add(key string) {
	c.ins++
	if c.ins == 0xffff {
		c.reset()
		c.ins = 0
	}

	h1, h2 := c.baseHashes(key)

	for i := 0; i < depth; i++ {
		pos := int(h1+uint32(i)*h2) & 0x1f
		shift := (pos & 1) * 4
		v := (c.nyb[(i*width+pos)/2] >> shift) & 0x0f
		if v < 15 {
			c.nyb[(i*width+pos)/2] += 1 << shift
		}
	}
}

func (c *CM4) baseHashes(key string) (uint32, uint32) {
	c.h.WriteString(key)
	h641 := c.h.Sum64()
	c.h.Reset()
	return uint32(h641), uint32(h641 >> 32)
}

func (c *CM4) Estimate(key string) uint8 {
	h1, h2 := c.baseHashes(key)

	var min byte = 255
	for i := 0; i < depth; i++ {
		pos := int(h1+uint32(i)*h2) & 0x1f
		off := i * width
		v := byte(c.nyb[(off+pos)/2]>>((pos&1)*4)) & 0x0f
		if v < min {
			min = v
		}
	}
	return min
}

func (c *CM4) reset() {
	for i := range c.nyb {
		// divide by two
		c.nyb[i] = (c.nyb[i] >> 1) & 0x77
	}
}
