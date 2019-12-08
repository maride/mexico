package interpreter

import (
	"fmt"
	"log"
)

func (t *Tape) DebugPrintTape() {
	log.Printf("Tape is currently %d cells big", len(t.cells))
	for i, v := range t.cells {
		fmt.Printf("Cell %d: %d (%c)\n", i, v, v)
	}
}

func (s * Stack) DebugPrintStack() {
	log.Printf("Stack is currently %d entries big", len(s.values))
	for i, v := range s.values {
		fmt.Printf("Stack row %d: %d (%c)\n", i, v, v)
	}
}
