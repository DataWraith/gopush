package gopush_test

import (
	"fmt"

	"github.com/DataWraith/gopush"
)

func Example() {
	// Use the default options
	options := gopush.DefaultOptions

	// Instantiate a new interpreter
	interpreter := gopush.NewInterpreter(options)

	// Create a new data type
	printStack := &gopush.Stack{
		Functions: make(map[string]func()),
	}

	// Add a function to the data type. This also demonstrates that
	// functions may have side effects outside of the interpreter when called.
	printStack.Functions["hello"] = func() {
		fmt.Println("hello")
	}

	// Register the new data type. The first statement adds the functions of
	// the type to the list of allowed functions, the second statement makes
	// the new type usable by the interpreter.
	interpreter.Options.RegisterStack("print", printStack)
	interpreter.RegisterStack("print", printStack)

	// Run the interpreter
	interpreter.Run("PRINT.HELLO")

	// Output: hello
}
