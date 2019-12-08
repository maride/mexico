package interpreter

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	Stack Stack
	Tape Tape
}

// Runs the given command.
// Returns the next line (comparable to the 'Program Counter') to be executed, but just if doJump is true.
// May also return an error. It's advised to stop the execution of further commands if this command throws an error.
func (m *Machine) RunCommand(cmd string) (jumpLine int, doJump bool, execErr error) {
	// Let's check which command we are told to run.
	if cmd == "left" {
		// Moves the tape head one cell to the left
		m.Tape.MoveLeft()
	} else if cmd == "right" {
		// Moves the tape head one cell to the right
		m.Tape.MoveRight()
	} else if cmd == "pusht" {
		// Reads the current cell value and pushes it on top of the stack
		m.Stack.Push(m.Tape.Get())
	} else if strings.HasPrefix(cmd, "push ") {
		// Pushes the value n to the stack

		// Cut "push " away and trim
		strVal := strings.Trim(cmd[5:], " ")

		// Convert to integer
		intVal, atoiErr := strconv.Atoi(strVal)
		if atoiErr != nil {
			// Conversion failed.
			execErr = errors.New(fmt.Sprintf("Tried to push non-integer value '%s' to the stack. %s", strVal, atoiErr.Error()))
			return
		} else {
			// Conversion successful, push constant
			m.Stack.Push(intVal)
		}
	} else if cmd == "pop" {
		// Pops top stack value to the current cell
		m.Tape.Set(m.Stack.Pop())
	} else if cmd == "dup" {
		// Duplicates the topmost stack value
		val := m.Stack.Pop()
		m.Stack.Push(val)
		m.Stack.Push(val)
	} else if cmd == "del" {
		// Deletes the topmost stack value, ignoring its value
		m.Stack.Pop()
	} else if cmd == "eq" {
		// Checks if stack[0] == stack[1]. Pushes 1 to the stack if equal, 0 otherwise
		stack0 := m.Stack.Pop()
		stack1 := m.Stack.Pop()

		if stack0 == stack1 {
			m.Stack.Push(1)
		} else {
			m.Stack.Push(0)
		}
	} else if cmd == "not" {
		// Inverses stack[0]
		stack0 := m.Stack.Pop()

		if stack0 == 0 {
			m.Stack.Push(1)
		} else if stack0 == 1 {
			m.Stack.Push(0)
		} else {
			// Not a binary number, not going to inverse it.
			execErr = errors.New(fmt.Sprintf("Tried to inverse non-binary integer value '%d'", stack0))
			return
		}
	} else if cmd == "gt" {
		// Checks if stack[0] > stack[1]. Pushes 1 to the stack if greater, 0 otherwise
		stack0 := m.Stack.Pop()
		stack1 := m.Stack.Pop()

		if stack0 > stack1 {
			m.Stack.Push(1)
		} else {
			m.Stack.Push(0)
		}
	} else if cmd == "lt" {
		// Checks if stack[0] < stack[1]. Pushes 1 to the stack if greater, 0 otherwise
		stack0 := m.Stack.Pop()
		stack1 := m.Stack.Pop()

		if stack0 < stack1 {
			m.Stack.Push(1)
		} else {
			m.Stack.Push(0)
		}
	} else if cmd == "add" {
		// Calculates stack[0] + stack[1], and pushes the result to the stack
		m.Stack.Push(m.Stack.Pop() + m.Stack.Pop())
	} else if cmd == "sub" {
		// Calculates stack[0] - stack[1], and pushes the result to the stack
		m.Stack.Push(m.Stack.Pop() - m.Stack.Pop())
	} else if cmd == "mult" {
		// Calculates stack[0] * stack[1], and pushes the result to the stack
		m.Stack.Push(m.Stack.Pop() * m.Stack.Pop())
	} else if cmd == "div" {
		// Calculates stack[0] / stack[1], and pushes the result to the stack
		m.Stack.Push(m.Stack.Pop() / m.Stack.Pop())
	} else if cmd == "mod" {
		// Calculates stack[0] % stack[1], and pushes the result to the stack
		m.Stack.Push(m.Stack.Pop() % m.Stack.Pop())
	} else if cmd == "read" {
		// Reads a character from the user, and pushes its char value to the stack
		var readChar []byte

		// Read character(s)
		numChars, readErr := os.Stdin.Read(readChar)
		if readErr != nil {
			// Failed to read from stdin.
			execErr = errors.New(fmt.Sprintf("Failed to read character from stdin: %s", readErr.Error()))
			return
		}

		if numChars < 1 || readChar == nil {
			// Didn't even read a single character.
			execErr = errors.New("Failed to read character from stdin.")
			return
		}

		// Push character to stack
		m.Stack.Push(int(readChar[0]))

		// Check if user typed more than asked for
		if numChars > 1 {
			// Ugh, what a spammer
			execErr = errors.New("Read more than one character - ignoring all but the first character.")
			return
		}
	} else if cmd == "print" {
		// Prints stack[0] as a character
		val := m.Stack.Pop()
		fmt.Printf("%q (%d)\n", rune(val), val)
	} else if cmd == "jmp" {
		// Jumps to the line number specified by stack[0]
		jumpLine = m.Stack.Pop()
		doJump = true
	} else if cmd == "jmpc" {
		// Jumps to the line number specified by stack[0], if stack[1] is not 0.
		jumpLine = m.Stack.Pop()
		doJump = m.Stack.Pop() != 0
	} else {
		// ... no such command.
		execErr = errors.New(fmt.Sprintf("Command not found: %s", cmd))
	}

	return
}
