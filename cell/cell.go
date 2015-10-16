// The cell package manages the state for the Conway's Game of Life cellular
// automaton.
package cell

import (
	"fmt"
	"image"
	"math/rand"
)

const (
	pixelBytes  = 4
	redMask     = 0xff0000
	redShift    = 16
	greenMask   = 0xff00
	greenShift  = 8
	greenOffset = 1
	greenColor  = 0x8000
	blueMask    = 0xff
	blueOffset  = 2
	alphaOffset = 3
	fullAlpha   = 0xff
)

// Field represents the entire state of a 2D cellular automaton.
type Field struct {
	// State holds the state information. The state for the cell at (x, y)
	// is stored at State[Width * y + x].
	State []uint32
	// Width is one dimension of the 2D state. The state is of size Width
	// times Height.
	Width int
	// Height is one dimension of the 2D state. The state is of size Width
	// times Height.
	Height int
}

// NewField creates an empty 2D cellular automaton state with the given (w, h)
// dimensions.
func NewField(w, h int) *Field {
	validateDimensions(w, h)
	return &Field{State: make([]uint32, w*h), Width: w, Height: h}
}

// RandomField creates a random 2D cellular automaton state with the given
// dimensions (w, h).
//
// It uses the lower 24 bits of state to make it easy to visualize using with
// an image (8 bits for each red, green, and blue).
func RandomField(w, h int) *Field {
	f := NewField(w, h)
	for i, _ := range f.State {
		f.State[i] = uint32(rand.Intn(1 << 24))
	}
	return f
}

func validateDimensions(w, h int) {
	if w < 0 {
		panic(fmt.Sprintf("Expected non-negative width, got %d", w))
	}
	if h < 0 {
		panic(fmt.Sprintf("Expected non-negative height, got %d", h))
	}
}

func validateField(f *Field) {
	validateDimensions(f.Width, f.Height)
	if len(f.State) != f.Width*f.Height {
		panic(fmt.Sprintf(
			"Expected field with State = %d = %d * %d, got %d",
			f.Width*f.Height,
			f.Width,
			f.Height,
			len(f.State)))
	}
}

// ToImage converts the lower 24 bits of a Field's state into an image.
func ToImage(f *Field) *image.NRGBA {
	validateField(f)
	img := image.NewNRGBA(image.Rect(0, 0, f.Width, f.Height))
	for i, s := range f.State {
		// Pull out the RGB values. I ignore the top alpha bit to make
		// visualizing the game of life easier. It would be almost
		// impossible to percieve differences in alpha values. They
		// will be 0 or close to 0 most of the time, which would hide
		// anything going on in the color channels.
		img.Pix[i*pixelBytes] = uint8((s & redMask) >> redShift)
		img.Pix[i*pixelBytes+greenOffset] = uint8((s & greenMask) >> greenShift)
		img.Pix[i*pixelBytes+blueOffset] = uint8(s & blueMask)
		img.Pix[i*pixelBytes+alphaOffset] = fullAlpha
	}
	return img
}