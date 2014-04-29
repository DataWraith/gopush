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
	interpreter.RegisterStack(stackName, printStack);

Finally, you can run the interpreter to execute a given program:

	program := "PRINT.HELLO"
	interpreter.Run(program)

*/
package gopush
