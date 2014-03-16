package gopush

func NewFloatStack(options Options) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["+"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f1 + f2)
	}

	s.Functions["*"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f1 * f2)
	}

	return s
}
