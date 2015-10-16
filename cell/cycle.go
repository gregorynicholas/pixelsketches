package cell

import (
	"github.com/bradfitz/iter"
	"time"
)

const threeMask = 0x49149249

// Run creates a channel which generates a new state for every message it
// recieves from tick, using the starting state c.
func Run(c *Field, tick <-chan time.Time) <-chan *Field {
	validateField(c)
	out := make(chan *Field)
	go func(c *Field, out chan<- *Field) {
		defer close(out)
		for _ = range tick {
			p := calc(c, plus)
			is2 := calc(p, two)
			is3 := calc(p, three)
			c = nextState(c, is2, is3)
			out <- c
		}
	}(c, out)
	return out
}

func nextState(c, is2, is3 *Field) *Field {
	out := NewField(c.Width, c.Height)
	for y := range iter.N(c.Height) {
		for x := range iter.N(c.Width) {
			i := y*c.Width + x
			// Live cells with 2 or 3 neighbors live.
			out.State[i] = ((is2.State[i] | is3.State[i]) & c.State[i]) |
				// Dead cells with exactly 3 neighbors live.
				(is3.State[i] &^ c.State[i])
		}
	}
	return out
}

type ceval func(c *Field, x, y int) uint32

// calc evaluates the function for each cell.
func calc(c *Field, f ceval) *Field {
	out := NewField(c.Width, c.Height)
	for y := range iter.N(c.Height) {
		for x := range iter.N(c.Width) {
			out.State[c.Width*y+x] = f(c, x, y)
		}
	}
	return out
}

func plus(c *Field, x, y int) uint32 {
	w := c.Width
	h := c.Height
	ym := fit(y-1, h)
	yp := fit(y+1, h)
	xm := fit(x-1, w)
	xp := fit(x+1, w)
	return c.State[w*yp+xm] + c.State[w*yp+x] + c.State[w*yp+xp] +
		c.State[w*y+xm] + c.State[w*y+xp] +
		c.State[w*ym+xm] + c.State[w*ym+x] + c.State[w*ym+xp]
}

// two calculates some bitwise operations which checks if the sum of neighbors
// is 2 for three-bit sections of the sum.
func two(c *Field, x, y int) uint32 {
	p := c.State[c.Width*y+x]
	// No bits should be in the 4's place.
	return (threeMask &^ (p >> 2)) & (
	// There should be a bit set in the 2's place.
	(threeMask & (p >> 1)) &
		// No bits should be set in the 1's place.
		(threeMask &^ p))
}

// three calculates some bitwise operations which checks if the sum of neighbors
// is 3 for three-bit sections of the sum.
func three(c *Field, x, y int) uint32 {
	p := c.State[c.Width*y+x]
	// No bits should be in the 4's place.
	return (threeMask &^ (p >> 2)) & (
	// There should be a bit set in the 2's place.
	(threeMask & (p >> 1)) &
		// There should be a bit set in the 1's place.
		(threeMask & p))
}

func fit(i, max int) int {
	i = i % max
	if i < 0 {
		return i + max
	}
	return i
}
