package gopush

// Stack represents a data type in the Push language. It contains the actual
// stack of values of that data type and a map of functions that pertain to that
// data type.
type Stack struct {
	Stack     []interface{}
	Functions map[string]func()
}

// Peek returns the topmost item on the stack, or an empty struct if the stack
// is empty.
func (s Stack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return struct{}{}
	}

	return s.Stack[len(s.Stack)-1]
}

// Push pushes a new element onto the stack.
func (s *Stack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit)
}

// Pop pops an element off the stack. It returns an empty struct if the stack is
// empty.
func (s *Stack) Pop() (item interface{}) {
	if len(s.Stack) == 0 {
		return struct{}{}
	}

	item = s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]

	return item
}

// Len returns the number of items on the stack.
func (s Stack) Len() int64 {
	return int64(len(s.Stack))
}

// Dup duplicates the item on top of the stack.
func (s *Stack) Dup() {
	if len(s.Stack) == 0 {
		return
	}

	s.Push(s.Peek())
}

// Swap swaps the top two items on the stack.
func (s *Stack) Swap() {
	if len(s.Stack) < 2 {
		return
	}

	i1 := s.Pop()
	i2 := s.Pop()
	s.Push(i1)
	s.Push(i2)
}

// Flush empties the stack
func (s *Stack) Flush() {
	s.Stack = nil
}

// Rot rotates the top three stack items by pulling out the third item and
// pushing it on top.
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

// Shove inserts an item deep into the stack, at index idx.
func (s *Stack) Shove(item interface{}, idx int64) {
	index := int64(len(s.Stack)-1) - idx
	if index < 0 {
		index = 0
	} else if index > int64(len(s.Stack)) {
		index = int64(len(s.Stack))
	}

	s.Stack = append(s.Stack[:index], append([]interface{}{item}, s.Stack[index:]...)...)
}

// Yank pulls out an item deep in the stack, at index idx, and puts it on top of
// the stack.
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

// YankDup copies an item deep in the stack, at index ids, and puts the copy on
// top of the stack.
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
