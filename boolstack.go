package gopush

import "fmt"

func NewBooleanStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("boolean", 2) {
			return
		}

		b1 := interpreter.Stacks["boolean"].Pop().(bool)
		b2 := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(b1 == b2)
	}

	s.Functions["and"] = func() {
		if !interpreter.stackOK("boolean", 2) {
			return
		}

		b1 := interpreter.Stacks["boolean"].Pop().(bool)
		b2 := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(b1 && b2)
	}

	s.Functions["define"] = func() {
		if !interpreter.stackOK("name", 1) || !interpreter.stackOK("boolean", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		b := interpreter.Stacks["boolean"].Pop().(bool)

		interpreter.Definitions[n] = Code{Length: 1, Literal: fmt.Sprint(b)}
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["boolean"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["boolean"].Flush()
	}

	s.Functions["fromfloat"] = func() {
		if !interpreter.stackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["boolean"].Push(f != 0)
	}

	s.Functions["frominteger"] = func() {
		if !interpreter.stackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i != 0)
	}

	s.Functions["not"] = func() {
		if !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(!b)
	}

	s.Functions["or"] = func() {
		if !interpreter.stackOK("boolean", 2) {
			return
		}

		b1 := interpreter.Stacks["boolean"].Pop().(bool)
		b2 := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(b1 || b2)
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["boolean"].Pop()
	}

	s.Functions["rand"] = func() {
		interpreter.Stacks["boolean"].Push(interpreter.Rand.Float64() < 0.5)
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["boolean"].Rot()
	}

	s.Functions["shove"] = func() {
		if !interpreter.stackOK("boolean", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Shove(b, i)
	}

	s.Functions["stackdepth"] = func() {
		if !interpreter.stackOK("integer", 0) {
			return
		}

		interpreter.Stacks["integer"].Push(interpreter.Stacks["boolean"].Len())
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["boolean"].Swap()
	}

	s.Functions["yank"] = func() {
		if !interpreter.stackOK("integer", 1) || !interpreter.stackOK("boolean", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Yank(i)
	}

	s.Functions["yankdup"] = func() {
		if !interpreter.stackOK("integer", 1) || !interpreter.stackOK("boolean", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].YankDup(i)
	}

	return s
}
