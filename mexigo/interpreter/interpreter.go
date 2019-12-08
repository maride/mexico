package interpreter

import (
	"fmt"
	"github.com/pkg/errors"
)

type Interpreter struct {
	machine Machine
	program []Codeline
	programCounter int
	programPointer int
}

// Feeds a new mexico machine with given code.
func Run(commands []Codeline) error {
	var i Interpreter

	// Set the given code as commands for the interpreter
	setCmdErr := i.SetCommands(commands)
	if setCmdErr != nil {
		return setCmdErr
	}

	// Let's run this program :)
	return i.Run()
}

// Sets the given array as new program for the interpreter
func (i *Interpreter) SetCommands(commands []Codeline) error {
	i.program = commands
	i.programCounter = 0
	return i.GoToNextCommand()
}

// Searches for the next command, starting from the current value of the programCounter.
// This may sound odd, because in most other architectures, this is just programCounter++, and there would be no need
// for a function like this. However, mexico has a BASIC-style program line numbering, means we need to search for the
// next line number containing code, because there may be one or more empty lines between the current and the next line.
// This is exactly what GoToNextCommand() does.
// If there is no next command, most likely because we reached the end of the program, an error is thrown.
func (i *Interpreter) GoToNextCommand() error {
	// Iterate over all lines to find the first which has a greater line number than the current one
	for index, line := range i.program {
		if line.Linenumber >= i.programCounter {
			// Found, set and return
			i.programCounter = line.Linenumber
			i.programPointer = index
			return nil
		}
	}

	// No next command found. Throw error.
	return errors.New(fmt.Sprintf("Found no commands after line %d. Stopping.", i.programCounter))
}

// Runs the commands, unless an error is encountered, then it doesn't run the commands.
func (i *Interpreter) Run() error {
	defer i.machine.Tape.DebugPrintTape()
	defer i.machine.Stack.DebugPrintStack()

	for {
		// Get current command
		cmd := i.program[i.programPointer]

		// Run command in the machine
		jumpLine, doJump, runErr := i.machine.RunCommand(cmd.Code)
		if runErr != nil {
			// Encountered an error during runtime, stop execution
			return runErr
		}

		// Check if we should jump anywhere else than to the next code line
		if doJump {
			// Yes, do it then.
			i.programCounter = jumpLine
		} else {
			// We are not asked to jump anywhere, move on to the next code line then.
			i.programCounter++
		}

		// Skip empty lines if there are any.
		nextCmdErr := i.GoToNextCommand()
		if nextCmdErr != nil {
			// Encountered error finding the next command, throw it
			return nextCmdErr
		}
	}
}
