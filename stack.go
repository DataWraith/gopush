package gopush

type Stack struct {
	Stack     []interface{}
	Functions map[string]func()
}

func (s Stack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return struct{}{}
	}

	return s.Stack[len(s.Stack)-1]
}

func (s *Stack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit)
}

func (s *Stack) Pop() (item interface{}) {
	if len(s.Stack) == 0 {
		return struct{}{}
	}

	item = s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]

	return item
}

func (s Stack) Len() int64 {
	return int64(len(s.Stack))
}

func (s *Stack) Dup() {
	if len(s.Stack) == 0 {
		return
	}

	s.Push(s.Peek())
}

func (s *Stack) Swap() {
	if len(s.Stack) < 2 {
		return
	}

	i1 := s.Pop()
	i2 := s.Pop()
	s.Push(i1)
	s.Push(i2)
}

func (s *Stack) Flush() {
	s.Stack = nil
}

func (s *Stack) Rot() {
	if len(s.Stack) < 3 {
		return
	}

	i1 := s.Pop()
	i2 := s.Pop()
	i3 := s.Pop()

	s.Push(i2)
	s.Push(i1)
	s.Push(i3)
}

func (s *Stack) Shove(item interface{}, idx int64) {
	index := int64(len(s.Stack)-1) - idx
	if index < 0 {
		index = 0
	} else if index > int64(len(s.Stack)) {
		index = int64(len(s.Stack))
	}

	s.Stack = append(s.Stack[:index], append([]interface{}{item}, s.Stack[index:]...)...)
}

func (s *Stack) Yank(idx int64) {
	if len(s.Stack) == 0 {
		return
	}

	index := int64(len(s.Stack)-1) - idx
	if index < 0 {
		index = 0
	} else if index > int64(len(s.Stack)-1) {
		index = int64(len(s.Stack) - 1)
	}

	item := s.Stack[index]
	s.Stack = append(s.Stack[:index], s.Stack[index+1:]...)
	s.Stack = append(s.Stack, item)
}

func (s *Stack) YankDup(idx int64) {
	if len(s.Stack) == 0 {
		return
	}

	index := int64(len(s.Stack)-1) - idx
	if index < 0 {
		index = 0
	} else if index > int64(len(s.Stack)-1) {
		index = int64(len(s.Stack) - 1)
	}

	item := s.Stack[index]
	s.Stack = append(s.Stack, item)
}
