package gopush

import (
	"strconv"
	"strings"
)

type Stack interface {
	Peek() interface{}
	Push(interface{})
}

type IntStack struct {
	Stack []int64
}

func (s IntStack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return 0
	}

	return s.Stack[len(s.Stack)-1]
}

func (s *IntStack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit.(int64))
}

type FloatStack struct {
	Stack []float64
}

func (s FloatStack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return 0.0
	}

	return s.Stack[len(s.Stack)-1]
}

func (s *FloatStack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit.(float64))
}

type Options struct {
}

type Interpreter struct {
	Stacks map[string]Stack
}

var DefaultOptions = Options{}

func NewInterpreter(options Options) *Interpreter {
	interpreter := &Interpreter{
		Stacks: make(map[string]Stack),
	}
	interpreter.Stacks["integer"] = new(IntStack)
	interpreter.Stacks["float"] = new(FloatStack)
	return interpreter
}

func (i *Interpreter) Run(program string) {
	instructions := strings.Split(program, " ")
	for _, instr := range instructions {
		// Parse Integer literal
		if intlit, err := strconv.ParseInt(instr, 10, 64); err == nil {
			i.Stacks["integer"].Push(intlit)
			continue
		}

		// Parse Float literal
		if floatlit, err := strconv.ParseFloat(instr, 64); err == nil {
			i.Stacks["float"].Push(floatlit)
			continue
		}
	}
}
