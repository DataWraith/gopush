package gopush

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Stack struct {
	Stack []interface{}
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

func (s Stack) Len() int {
	return len(s.Stack)
}

type Options struct {
}

type Interpreter struct {
	Stacks map[string]*Stack
}

var DefaultOptions = Options{}

func NewInterpreter(options Options) *Interpreter {
	interpreter := &Interpreter{
		Stacks: make(map[string]*Stack),
	}

	interpreter.Stacks["integer"] = new(Stack)
	interpreter.Stacks["float"] = new(Stack)
	interpreter.Stacks["exec"] = new(Stack)
	interpreter.Stacks["boolean"] = new(Stack)

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

		if boollit, err := strconv.ParseBool(item); err == nil {
			i.Stacks["boolean"].Push(boollit)
			continue
		}

		// Try to parse the item on top of the exec stack as instruction
		// TODO
	}

	return nil
}
