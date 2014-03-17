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

	s.Functions["do"] = func() {
		if !interpreter.stackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)
		interpreter.runCode(c)
		interpreter.Stacks["code"].Pop()
	}

	s.Functions["do*range"] = func() {
		if !interpreter.stackOK("code", 1) || !interpreter.stackOK("integer", 2) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)
		dst := interpreter.Stacks["integer"].Pop().(int64)
		cur := interpreter.Stacks["integer"].Pop().(int64)

		if cur == dst {
			interpreter.Stacks["integer"].Push(cur)
			interpreter.Stacks["exec"].Push(c)
		} else {
			interpreter.Stacks["integer"].Push(cur)

			if dst < cur {
				cur--
			} else {
				cur++
			}

			interpreter.Stacks["code"].Push(c)
			interpreter.Stacks["exec"].Push(c)
			interpreter.Stacks["exec"].Push(Code{Length: 1, Literal: "CODE.DO*RANGE"})
			interpreter.Stacks["integer"].Push(cur)
			interpreter.Stacks["integer"].Push(dst)
		}

	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["code"].Dup()
	}

	s.Functions["if"] = func() {
		if !interpreter.stackOK("code", 2) || !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if b {
			interpreter.Stacks["exec"].Push(c2)
		} else {
			interpreter.Stacks["exec"].Push(c1)
		}
	}

	s.Functions["quote"] = func() {
		if !interpreter.stackOK("exec", 1) {
			return
		}

		c := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["code"].Push(c)
	}

	return s
}
