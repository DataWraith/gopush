package gopush

import "reflect"

func newCodeStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("code", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if reflect.DeepEqual(c1, c2) {
			interpreter.Stacks["boolean"].Push(true)
		} else {
			interpreter.Stacks["boolean"].Push(false)
		}
	}

	s.Functions["append"] = func() {
		if !interpreter.stackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if c1.Literal != "" {
			c1 = Code{Length: c1.Length, List: []Code{c1}}
		}

		if c2.Literal != "" {
			c2 = Code{Length: c2.Length, List: []Code{c2}}
		}

		combined := Code{Length: c1.Length + c2.Length, List: append(c2.List, c1.List...)}

		if combined.Length <= interpreter.Options.MaxPointsInProgram {
			interpreter.Stacks["code"].Push(combined)
		}
	}

	s.Functions["atom"] = func() {
		if !interpreter.stackOK("code", 1) || !interpreter.stackOK("boolean", 0) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if c.Literal != "" {
			interpreter.Stacks["boolean"].Push(true)
		} else {
			interpreter.Stacks["boolean"].Push(false)
		}
	}

	s.Functions["car"] = func() {
		if !interpreter.stackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if len(c.List) == 0 {
			return
		}

		interpreter.Stacks["code"].Push(c.List[0])
	}

	s.Functions["cdr"] = func() {
		if !interpreter.stackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if len(c.List) == 0 {
			interpreter.Stacks["code"].Push(Code{})
		} else {
			cdr := Code{
				Length: c.Length - c.List[0].Length,
				List:   c.List[1:],
			}
			interpreter.Stacks["code"].Push(cdr)
		}
	}

	s.Functions["cons"] = func() {
		if !interpreter.stackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if c1.Literal != "" {
			c1 = Code{Length: 1, List: []Code{c1}}
		}

		if c2.Literal != "" {
			c2 = Code{Length: 1, List: []Code{c2}}
		}

		c := Code{
			Length: c1.Length + c2.Length,
			List:   append(c2.List, c1.List...),
		}

		interpreter.Stacks["code"].Push(c)
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
		interpreter.Stacks["code"].Flush()
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
		// Does nothing
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
		if !interpreter.stackOK("integer", 1) {
			return
		}

		maxPoints := interpreter.Stacks["integer"].Pop().(int64)
		if maxPoints < 0 {
			maxPoints *= -1
		}

		if maxPoints > interpreter.Options.MaxPointsInRandomExpression {
			maxPoints = interpreter.Options.MaxPointsInRandomExpression
		}

		c := interpreter.RandomCode(maxPoints)
		interpreter.Stacks["code"].Push(c)
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
