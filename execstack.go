package gopush

func newExecStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["do*range"] = func() {
		if !interpreter.stackOK("exec", 1) || !interpreter.stackOK("integer", 2) {
			return
		}

		c := interpreter.Stacks["exec"].Pop().(Code)
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

			interpreter.Stacks["integer"].Push(cur)
			interpreter.Stacks["integer"].Push(dst)
			interpreter.Stacks["exec"].Push(c)
			interpreter.Stacks["exec"].Push(c)
			interpreter.Stacks["exec"].Push(Code{Length: 1, Literal: "EXEC.DO*RANGE"})
		}
	}

	s.Functions["if"] = func() {
		if !interpreter.stackOK("exec", 2) || !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		c1 := interpreter.Stacks["exec"].Pop().(Code)
		c2 := interpreter.Stacks["exec"].Pop().(Code)

		if b {
			interpreter.Stacks["exec"].Push(c1)
		} else {
			interpreter.Stacks["exec"].Push(c2)
		}
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["exec"].Pop()
	}

	s.Functions["y"] = func() {
		if !interpreter.stackOK("exec", 1) {
			return
		}

		e := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["exec"].Push(Code{Length: 2, List: []Code{Code{Length: 1, Literal: "EXEC.Y"}, e}})
		interpreter.Stacks["exec"].Push(e)
	}

	return s
}
