package gopush

func NewCodeStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["define"] = func() {
		if !interpreter.stackOK("name", 1) || !interpreter.stackOK("code", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		c := interpreter.Stacks["code"].Pop().(Code)

		interpreter.Definitions[n] = Definition{Stack: "exec", Value: c}
	}

	s.Functions["quote"] = func() {
		c := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["code"].Push(c)
	}

	return s
}
