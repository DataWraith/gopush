package gopush

import "fmt"

func NewIntStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["%"] = func() {
		if !interpreter.stackOK("integer", 2) || interpreter.Stacks["integer"].Peek().(int64) == 0 {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)

		mod := i2 % i1
		if (i2 < 0 && i1 > 0) || (i2 > 0 && i1 < 0) {
			mod = i1 + mod
		}

		interpreter.Stacks["integer"].Push(mod)
	}

	s.Functions["*"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Push(i1 * i2)
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

	s.Functions["/"] = func() {
		if !interpreter.stackOK("integer", 2) || interpreter.Stacks["integer"].Peek().(int64) == 0 {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Push(i2 / i1)
	}

	s.Functions["<"] = func() {
		if !interpreter.stackOK("integer", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i2 < i1)
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("integer", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i2 == i1)
	}

	s.Functions[">"] = func() {
		if !interpreter.stackOK("integer", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i2 > i1)
	}

	s.Functions["define"] = func() {
		if !interpreter.stackOK("name", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		i := interpreter.Stacks["integer"].Pop().(int64)

		interpreter.Definitions[n] = Code{Length: 1, Literal: fmt.Sprint(i)}
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["integer"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["integer"].Flush()
	}

	s.Functions["fromboolean"] = func() {
		if !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		if b {
			interpreter.Stacks["integer"].Push(int64(1))
		} else {
			interpreter.Stacks["integer"].Push(int64(0))
		}
	}

	s.Functions["fromfloat"] = func() {
		if !interpreter.stackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		interpreter.Stacks["integer"].Push(int64(f))
	}

	s.Functions["max"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)

		if i1 > i2 {
			interpreter.Stacks["integer"].Push(i1)
		} else {
			interpreter.Stacks["integer"].Push(i2)
		}
	}

	s.Functions["min"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)

		if i1 < i2 {
			interpreter.Stacks["integer"].Push(i1)
		} else {
			interpreter.Stacks["integer"].Push(i2)
		}
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["integer"].Pop()
	}

	s.Functions["rand"] = func() {
		// TODO
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["integer"].Rot()
	}

	s.Functions["shove"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Peek().(int64)

		interpreter.Stacks["integer"].Shove(i2, i1)
		interpreter.Stacks["integer"].Pop()
	}

	s.Functions["stackdepth"] = func() {
		interpreter.Stacks["integer"].Push(interpreter.Stacks["integer"].Len())
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["integer"].Swap()
	}

	s.Functions["yank"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].Yank(idx)
	}

	s.Functions["yankdup"] = func() {
		if !interpreter.stackOK("integer", 2) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["integer"].YankDup(idx)
	}

	return s
}
