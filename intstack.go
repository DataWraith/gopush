package gopush

func NewIntStack(options Options) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["+"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() < 2 {
			return
		}

		i1 := stacks["integer"].Pop().(int64)
		i2 := stacks["integer"].Pop().(int64)
		stacks["integer"].Push(i1 + i2)
	}

	s.Functions["-"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() < 2 {
			return
		}

		i1 := stacks["integer"].Pop().(int64)
		i2 := stacks["integer"].Pop().(int64)
		stacks["integer"].Push(i2 - i1)
	}

	s.Functions["*"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() < 2 {
			return
		}

		i1 := stacks["integer"].Pop().(int64)
		i2 := stacks["integer"].Pop().(int64)
		stacks["integer"].Push(i1 * i2)
	}

	return s
}
