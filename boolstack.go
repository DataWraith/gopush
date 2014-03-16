package gopush

import "math/rand"

func NewBooleanStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		if interpreter.Stacks["boolean"].Len() < 2 {
			return
		}

		b1 := interpreter.Stacks["boolean"].Pop().(bool)
		b2 := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(b1 == b2)
	}

	s.Functions["and"] = func() {
		if interpreter.Stacks["boolean"].Len() < 2 {
			return
		}

		b1 := interpreter.Stacks["boolean"].Pop().(bool)
		b2 := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(b1 && b2)
	}

	s.Functions["define"] = func() {
		// TODO
		return
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["boolean"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["boolean"].Flush()
	}

	s.Functions["fromfloat"] = func() {
		if interpreter.Stacks["float"].Len() == 0 {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["boolean"].Push(f != 0)
	}

	s.Functions["frominteger"] = func() {
		if interpreter.Stacks["integer"].Len() == 0 {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i != 0)
	}

	s.Functions["not"] = func() {
		if interpreter.Stacks["boolean"].Len() == 0 {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["boolean"].Push(!b)
	}

	s.Functions["or"] = func() {
		if interpreter.Stacks["boolean"].Len() < 2 {
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
		interpreter.Stacks["boolean"].Push(rand.Float64() < 0.5)
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["boolean"].Rot()
	}

	s.Functions["shove"] = func() {
		if interpreter.Stacks["boolean"].Len() == 0 || interpreter.Stacks["integer"].Len() == 0 {
			return
		}

		interpreter.Stacks["boolean"].Shove(interpreter.Stacks["boolean"].Pop(), interpreter.Stacks["integer"].Pop().(int64))
	}

	s.Functions["stackdepth"] = func() {
		interpreter.Stacks["integer"].Push(interpreter.Stacks["boolean"].Len())
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["boolean"].Swap()
	}

	s.Functions["yank"] = func() {
		if interpreter.Stacks["integer"].Len() == 0 {
			return
		}

		interpreter.Stacks["boolean"].Yank(interpreter.Stacks["integer"].Pop().(int64))
	}

	s.Functions["yankdup"] = func() {
		if interpreter.Stacks["integer"].Len() == 0 {
			return
		}

		interpreter.Stacks["boolean"].YankDup(interpreter.Stacks["integer"].Pop().(int64))
	}

	return s
}
