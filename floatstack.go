package gopush

import "math"

func NewFloatStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["%"] = func() {
		// TODO
	}

	s.Functions["*"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(f1 * f2)
	}

	s.Functions["+"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(f1 + f2)
	}

	s.Functions["-"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(f2 - f1)
	}

	s.Functions["/"] = func() {
		if !interpreter.stackOK("float", 2) || interpreter.Stacks["float"].Peek().(float64) == 0 {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(f2 / f1)
	}

	s.Functions["<"] = func() {
		if !interpreter.stackOK("float", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["boolean"].Push(f2 < f1)
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("float", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["boolean"].Push(f1 == f2)
	}

	s.Functions[">"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["boolean"].Push(f2 > f1)
	}

	s.Functions["cos"] = func() {
		if !interpreter.stackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(math.Cos(f))
	}

	s.Functions["define"] = func() {
		// TODO
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["float"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["float"].Flush()
	}

	s.Functions["fromboolean"] = func() {
		if !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		if b {
			interpreter.Stacks["float"].Push(1.0)
		} else {
			interpreter.Stacks["float"].Push(0.0)
		}
	}

	s.Functions["frominteger"] = func() {
		if !interpreter.stackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["float"].Push(float64(i))
	}

	s.Functions["max"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(math.Max(f1, f2))
	}

	s.Functions["min"] = func() {
		if !interpreter.stackOK("float", 2) {
			return
		}

		f1 := interpreter.Stacks["float"].Pop().(float64)
		f2 := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(math.Min(f1, f2))
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["float"].Pop()
	}

	s.Functions["rand"] = func() {
		// TODO
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["float"].Rot()
	}

	s.Functions["shove"] = func() {
		if !interpreter.stackOK("float", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["float"].Shove(f, i)
	}

	s.Functions["sin"] = func() {
		if !interpreter.stackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(math.Sin(f))
	}

	s.Functions["stackdepth"] = func() {
		if !interpreter.stackOK("integer", 0) {
			return
		}

		interpreter.Stacks["integer"].Push(interpreter.Stacks["float"].Len())
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["float"].Swap()
	}

	s.Functions["tan"] = func() {
		if !interpreter.stackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["float"].Push(math.Tan(f))
	}

	s.Functions["yank"] = func() {
		if !interpreter.stackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["float"].Yank(i)
	}

	s.Functions["yankdup"] = func() {
		if !interpreter.stackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["float"].YankDup(i)
	}

	return s
}
