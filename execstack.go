package gopush

func NewExecStack(interpreter *Interpreter) *Stack {
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

	return s
}
