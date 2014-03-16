package gopush

func NewIntStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["+"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Push(i1 + i2)
	}

	s.Functions["-"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Push(i2 - i1)
	}

	s.Functions["*"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Push(i1 * i2)
	}

	s.Functions["<"] = func() {
		if !interpreter.stackOK("integer", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i2 < i1)
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["integer"].Dup()
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["integer"].Pop()
	}

	return s
}
