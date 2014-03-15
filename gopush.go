package gopush

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Instruction func(map[string]*Stack)

type Stack struct {
	Stack     []interface{}
	Functions map[string]Instruction
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

func NewBooleanStack(options Options) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["or"] = func(stacks map[string]*Stack) {
		if stacks["boolean"].Len() < 2 {
			return
		}

		b1 := stacks["boolean"].Pop().(bool)
		b2 := stacks["boolean"].Pop().(bool)
		stacks["boolean"].Push(b1 || b2)
	}

	return s
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

	interpreter.Stacks["integer"] = NewIntStack(options)
	interpreter.Stacks["float"] = NewFloatStack(options)
	interpreter.Stacks["exec"] = new(Stack)
	interpreter.Stacks["boolean"] = NewBooleanStack(options)

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

		if t == "" {
			continue
		}

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
			continue
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
		item = strings.ToLower(item)
		if strings.Contains(item, ".") {
			stack := item[:strings.Index(item, ".")]
			operation := item[strings.Index(item, ".")+1:]

			s, ok := i.Stacks[stack]
			if !ok {
				return errors.New(fmt.Sprintf("unkown stack: %v", stack))
			}

			f, ok := s.Functions[operation]
			if !ok {
				return errors.New(fmt.Sprintf("unknown instruction: %v.%v", stack, operation))
			}

			f(i.Stacks)
			continue
		}

		return errors.New(fmt.Sprintf("not an instruction: %q", item))
	}

	return nil
}
