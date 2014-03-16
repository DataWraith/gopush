package gopush

import "math/rand"

func NewBooleanStack(options Options) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() < 2 {
			return
		}

		b1 := stacks["boolean"].Pop().(bool)
		b2 := stacks["boolean"].Pop().(bool)
		stacks["boolean"].Push(b1 == b2)
	}

	s.Functions["and"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() < 2 {
			return
		}

		b1 := stacks["boolean"].Pop().(bool)
		b2 := stacks["boolean"].Pop().(bool)
		stacks["boolean"].Push(b1 && b2)
	}

	s.Functions["define"] = func(stacks map[string]*Stack) {
		// TODO
		return
	}

	s.Functions["dup"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Dup()
	}

	s.Functions["flush"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Flush()
	}

	s.Functions["fromfloat"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() == 0 {
			return
		}

		f := stacks["float"].Pop().(float64)
		stacks["boolean"].Push(f != 0)
	}

	s.Functions["frominteger"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		i := stacks["integer"].Pop().(int64)
		stacks["boolean"].Push(i != 0)
	}

	s.Functions["not"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() == 0 {
			return
		}

		b := stacks["boolean"].Pop().(bool)
		stacks["boolean"].Push(!b)
	}

	s.Functions["or"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() < 2 {
			return
		}

		b1 := stacks["boolean"].Pop().(bool)
		b2 := stacks["boolean"].Pop().(bool)
		stacks["boolean"].Push(b1 || b2)
	}

	s.Functions["pop"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Pop()
	}

	s.Functions["rand"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Push(rand.Float64() < 0.5)
	}

	s.Functions["rot"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Rot()
	}

	s.Functions["shove"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() == 0 || stacks["integer"].Len() == 0 {
			return
		}

		stacks["boolean"].Shove(stacks["boolean"].Pop(), stacks["integer"].Pop().(int64))
	}

	s.Functions["stackdepth"] = func(stacks map[string]*Stack) {
		stacks["integer"].Push(stacks["boolean"].Len())
	}

	s.Functions["swap"] = func(stacks map[string]*Stack) {
		stacks["boolean"].Swap()
	}

	s.Functions["yank"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		stacks["boolean"].Yank(stacks["integer"].Pop().(int64))
	}

	s.Functions["yankdup"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		stacks["boolean"].YankDup(stacks["integer"].Pop().(int64))
	}

	return s
}
