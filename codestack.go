package gopush

func newCodeStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		// TODO
	}

	s.Functions["append"] = func() {
		// TODO
	}

	s.Functions["atom"] = func() {
		// TODO
	}

	s.Functions["car"] = func() {
		// TODO
	}

	s.Functions["cdr"] = func() {
		// TODO
	}

	s.Functions["cons"] = func() {
		// TODO
	}

	s.Functions["container"] = func() {
		// TODO
	}

	s.Functions["contains"] = func() {
		// TODO
	}

	s.Functions["define"] = func() {
		if !interpreter.stackOK("name", 1) || !interpreter.stackOK("code", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		c := interpreter.Stacks["code"].Pop().(Code)

		interpreter.define(n, c)
	}

	s.Functions["definition"] = func() {
		// TODO
	}

	s.Functions["discrepancy"] = func() {
		// TODO
	}

	s.Functions["do"] = func() {
		if !interpreter.stackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		err := interpreter.runCode(c)
		if err != nil {
			panic(err)
		}

		interpreter.Stacks["code"].Pop()
	}

	s.Functions["do*"] = func() {
		// TODO
	}

	s.Functions["do*count"] = func() {
		// TODO
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

	s.Functions["do*times"] = func() {
		// TODO
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["code"].Dup()
	}

	s.Functions["extract"] = func() {
		// TODO
	}

	s.Functions["flush"] = func() {
		// TODO
	}

	s.Functions["fromboolean"] = func() {
		// TODO
	}

	s.Functions["fromfloat"] = func() {
		// TODO
	}

	s.Functions["frominteger"] = func() {
		// TODO
	}

	s.Functions["fromname"] = func() {
		// TODO
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

	s.Functions["insert"] = func() {
		// TODO
	}

	s.Functions["instructions"] = func() {
		// TODO
	}

	s.Functions["length"] = func() {
		// TODO
	}

	s.Functions["list"] = func() {
		// TODO
	}

	s.Functions["member"] = func() {
		// TODO
	}

	s.Functions["noop"] = func() {
		// TODO
	}

	s.Functions["nth"] = func() {
		// TODO
	}

	s.Functions["nthcdr"] = func() {
		// TODO
	}

	s.Functions["null"] = func() {
		// TODO
	}

	s.Functions["pop"] = func() {
		// TODO
	}

	s.Functions["position"] = func() {
		// TODO
	}

	s.Functions["quote"] = func() {
		if !interpreter.stackOK("exec", 1) {
			return
		}

		c := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["rand"] = func() {
		// TODO
	}

	s.Functions["rot"] = func() {
		// TODO
	}

	s.Functions["shove"] = func() {
		// TODO
	}

	s.Functions["size"] = func() {
		// TODO
	}

	s.Functions["stackdepth"] = func() {
		// TODO
	}

	s.Functions["subst"] = func() {
		// TODO
	}

	s.Functions["swap"] = func() {
		// TODO
	}

	s.Functions["yank"] = func() {
		// TODO
	}

	s.Functions["yankdup"] = func() {
		// TODO
	}

	return s
}
