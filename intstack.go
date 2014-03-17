package gopush

import "fmt"

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

	s.Functions["="] = func() {
		if !interpreter.stackOK("integer", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		i1 := interpreter.Stacks["integer"].Pop().(int64)
		i2 := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["boolean"].Push(i2 == i1)
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

	s.Functions["pop"] = func() {
		interpreter.Stacks["integer"].Pop()
	}

	return s
}
