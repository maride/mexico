package interpreter

type Tape struct {
	head uint
	cells []int
}

// Sets the head to point to the specified position
func (t *Tape) SetHead(pos uint) {
	t.GrowUpTo(pos)
	t.head = pos
}

// Resizes the tape to the given size.
// If the cell array is already bigger than the specified size, nothing happens.
// If the cell array is smaller than the specified size, it's resized to the given size and filled with the value 0.
func (t *Tape) GrowUpTo(size uint) {
	if uint(len(t.cells)) > size {
		// Tape is already bigger. Do nothing.
		return
	}

	// Append as many zeroes as required.
	for uint(len(t.cells)) <= size {
		t.cells = append(t.cells, 0)
	}
}

// Moves the tape head left
func (t *Tape) MoveLeft() {
	if t.head > 0 {
		t.SetHead(t.head - 1)
	}
}

// Moves the tape head left
func (t *Tape) MoveRight() {
	t.SetHead(t.head + 1)
}

// Returns the current cell value
func (t *Tape) Get() int {
	t.GrowUpTo(t.head)
	return t.cells[t.head]
}

// Sets the current cell to the new value
func (t *Tape) Set(newVal int) {
	t.GrowUpTo(t.head)
	t.cells[t.head] = newVal
}
