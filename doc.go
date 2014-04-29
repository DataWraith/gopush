/*

Package gopush provides an interpreter for the Push 3.0 programming language.

To create a new Interpreter, you must first specify a set of Options, or use the
provided DefaultOptions:

	options := gopush.DefaultOptions
	interpreter := gopush.NewInterpreter(options)

You can provide custom data types and associated behavior by implementing a new
Stack object:

	printStack := &gopush.Stack{
		Functions: make(map[string]func())
	}

	printStack.Functions["hello"] = func() {
		fmt.Println("hello")
	}

The keys of the Functions map *must* be lowercase. For more information on
stacks, take a look at the builtin stacks in the stack_* files.

After creating your new data type, you need to register it with Interpreter to
make it usable.

	stackName := "print"
	interpreter.Options.RegisterStack(stackName, printStack)
	interpreter.RegisterStack(stackName, printStack)

The first RegisterStack call adds all instructions from printStack to the list
of allowed instructions. The second RegisterStack call adds printStack and its
instructions to the Interpreter so they can be used. This double-registration is
not ideal, and I would welcome suggestions on how to remedy it without losing
the ability to mark certain instructions as disallowed.

Finally, you can run the interpreter to execute a given program:

	program := "PRINT.HELLO"
	interpreter.Run(program)

Alternatively you can parse the program into the Code representation and run
that:

	program := "PRINT.HELLO"
	c, err := ParseCode(program)
	if err != nil {
		// handle error
	}

	interpreter.RunCode(c)

*/
package gopush
