package gopush

import "github.com/cryptix/goremutake"

// newNameStack returns a new NAME stack
func newNameStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("name", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		n1 := interpreter.Stacks["name"].Pop().(string)
		n2 := interpreter.Stacks["name"].Pop().(string)
		interpreter.Stacks["boolean"].Push(n1 == n2)
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["name"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["name"].Flush()
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["name"].Pop()
	}

	s.Functions["quote"] = func() {
		interpreter.quoteNextName = true
	}

	s.Functions["rand"] = func() {
		randName := goremutake.Encode(interpreter.numNamesGenerated)
		interpreter.Stacks["name"].Push(randName)
		interpreter.numNamesGenerated++
	}

	s.Functions["randboundname"] = func() {
		// TODO
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["name"].Rot()
	}

	s.Functions["shove"] = func() {
		if !interpreter.stackOK("name", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		name := interpreter.Stacks["name"].Peek().(string)
		interpreter.Stacks["name"].Shove(name, idx)
		interpreter.Stacks["name"].Pop()
	}

	s.Functions["stackdepth"] = func() {
		if !interpreter.stackOK("integer", 0) {
			return
		}

		interpreter.Stacks["integer"].Push(interpreter.Stacks["name"].Len())
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["name"].Swap()
	}

	s.Functions["yank"] = func() {
		if !interpreter.stackOK("integer", 1) || !interpreter.stackOK("name", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["name"].Yank(idx)
	}

	s.Functions["yankdup"] = func() {
		if !interpreter.stackOK("integer", 1) || !interpreter.stackOK("name", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["name"].YankDup(idx)
	}

	return s
}
