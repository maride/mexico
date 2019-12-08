package interpreter

import "log"

type Stack struct {
	values []int
}

// Pushes the given value to the stack
func (s *Stack) Push(val int) {
	s.values = append(s.values, val)
}

// Pops the top element from the stack and return its value
func (s *Stack) Pop() int {
	// Check if the stack contains at least one element
	if len(s.values) > 0 {
		// There's at least one element, pop it: get value and delete element
		val := s.values[len(s.values)-1]
		s.values = s.values[:len(s.values)-1]
		return val
	}

	// Stack is empty, but we should pop... Damn.
	log.Panic("Tried to pop value from empty stack.")
	return 0
}
