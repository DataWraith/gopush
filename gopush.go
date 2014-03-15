package gopush

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Stack interface {
	Peek() interface{}
	Push(interface{})
	Pop() interface{}
	Len() int
}

type IntStack struct {
	Stack []int64
}

func (s IntStack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return int64(0)
	}

	return s.Stack[len(s.Stack)-1]
}

func (s *IntStack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit.(int64))
}

func (s *IntStack) Pop() (item interface{}) {
	if len(s.Stack) == 0 {
		return int64(0)
	}

	item = s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return item
}

func (s IntStack) Len() int {
	return len(s.Stack)
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

func (s *FloatStack) Pop() (item interface{}) {
	if len(s.Stack) == 0 {
		return 0.0
	}

	item = s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return item
}

func (s FloatStack) Len() int {
	return len(s.Stack)
}

type ExecStack struct {
	Stack []string
}

func (s ExecStack) Peek() interface{} {
	if len(s.Stack) == 0 {
		return ""
	}

	return s.Stack[len(s.Stack)-1]
}

func (s *ExecStack) Push(lit interface{}) {
	s.Stack = append(s.Stack, lit.(string))
}

func (s *ExecStack) Pop() (item interface{}) {
	if len(s.Stack) == 0 {
		return ""
	}

	item = s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return item
}

func (s ExecStack) Len() int {
	return len(s.Stack)
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
	interpreter.Stacks["exec"] = new(ExecStack)
	return interpreter
}

func ignoreWhiteSpace(program string) string {
	for i, r := range program {
		if !unicode.IsSpace(r) {
			return program[i:]
		}
	}
	return ""
}

func getToken(program string) (token, remainder string) {
	for i, r := range program {
		if unicode.IsSpace(r) {
			return program[:i], program[i:]
		}
	}
	return program, ""
}

func getToParen(program string) (subprogram, remainder string, err error) {
	parenBalance := 1
	for i, r := range program {
		switch r {
		case '(':
			parenBalance++
		case ')':
			parenBalance--
		}

		if parenBalance == 0 {
			return program[:i], program[i+1:], nil
		}
	}
	return "", "", errors.New("unmatched parentheses")
}

func splitProgram(program string) (result []string, err error) {
	var p, t string

	p = program
	for len(p) > 0 {
		p = ignoreWhiteSpace(p)
		t, p = getToken(p)

		if t == "(" {
			t, p, err = getToParen(p)
			if err != nil {
				return []string{}, err
			}
		}

		result = append(result, t)
	}

	return result, nil
}

func (i *Interpreter) Run(program string) (err error) {
	i.Stacks["exec"].Push(strings.TrimSpace(program))

	for i.Stacks["exec"].Len() > 0 {
		item := i.Stacks["exec"].Pop().(string)

		// If the item on top of the exec stack is a list, push it in
		// reverse order
		if strings.Contains(item, " ") {
			p, err := splitProgram(item)
			if err != nil {
				return err
			}
			for j := len(p) - 1; j >= 0; j-- {
				i.Stacks["exec"].Push(p[j])
			}
		}

		// Try to parse the item on top of the exec stack as a literal
		if intlit, err := strconv.ParseInt(item, 10, 64); err == nil {
			i.Stacks["integer"].Push(intlit)
			continue
		}

		if floatlit, err := strconv.ParseFloat(item, 64); err == nil {
			i.Stacks["float"].Push(floatlit)
			continue
		}

		// Try to parse the item on top of the exec stack as instruction
		// TODO
	}

	return nil
}
