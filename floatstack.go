package gopush

import "math"

func NewFloatStack(options Options) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["%"] = func(stacks map[string]*Stack) {
		// TODO
	}

	s.Functions["*"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f1 * f2)
	}

	s.Functions["+"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f1 + f2)
	}

	s.Functions["-"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f2 - f1)
	}

	s.Functions["/"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 || stacks["float"].Peek().(float64) == 0 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(f2 / f1)
	}

	s.Functions["<"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["boolean"].Push(f2 < f1)
	}

	s.Functions["="] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["boolean"].Push(f1 == f2)
	}

	s.Functions[">"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["boolean"].Push(f2 > f1)
	}

	s.Functions["cos"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() == 0 {
			return
		}

		f := stacks["float"].Pop().(float64)
		stacks["float"].Push(math.Cos(f))
	}

	s.Functions["define"] = func(stacks map[string]*Stack) {
		// TODO
	}

	s.Functions["dup"] = func(stacks map[string]*Stack) {
		stacks["float"].Dup()
	}

	s.Functions["flush"] = func(stacks map[string]*Stack) {
		stacks["float"].Flush()
	}

	s.Functions["fromboolean"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() == 0 {
			return
		}

		b := stacks["boolean"].Pop().(bool)
		if b {
			stacks["float"].Push(1.0)
		} else {
			stacks["float"].Push(0.0)
		}
	}

	s.Functions["frominteger"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		i := stacks["integer"].Pop().(int64)
		stacks["float"].Push(float64(i))
	}

	s.Functions["max"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(math.Max(f1, f2))
	}

	s.Functions["min"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() < 2 {
			return
		}

		f1 := stacks["float"].Pop().(float64)
		f2 := stacks["float"].Pop().(float64)
		stacks["float"].Push(math.Min(f1, f2))
	}

	s.Functions["pop"] = func(stacks map[string]*Stack) {
		stacks["float"].Pop()
	}

	s.Functions["rand"] = func(stacks map[string]*Stack) {
		// TODO
	}

	s.Functions["rot"] = func(stacks map[string]*Stack) {
		stacks["float"].Rot()
	}

	s.Functions["shove"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() == 0 || stacks["integer"].Len() == 0 {
			return
		}

		f := stacks["float"].Pop().(float64)
		i := stacks["integer"].Pop().(int64)
		stacks["float"].Shove(f, i)
	}

	s.Functions["sin"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() == 0 {
			return
		}

		f := stacks["float"].Pop().(float64)
		stacks["float"].Push(math.Sin(f))
	}

	s.Functions["stackdepth"] = func(stacks map[string]*Stack) {
		stacks["integer"].Push(stacks["float"].Len())
	}

	s.Functions["swap"] = func(stacks map[string]*Stack) {
		stacks["float"].Swap()
	}

	s.Functions["tan"] = func(stacks map[string]*Stack) {
		if stacks["float"].Len() == 0 {
			return
		}

		f := stacks["float"].Pop().(float64)
		stacks["float"].Push(math.Tan(f))
	}

	s.Functions["yank"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		i := stacks["integer"].Pop().(int64)
		stacks["float"].Yank(i)
	}

	s.Functions["yankdup"] = func(stacks map[string]*Stack) {
		if stacks["integer"].Len() == 0 {
			return
		}

		i := stacks["integer"].Pop().(int64)
		stacks["float"].YankDup(i)
	}

	return s
}
