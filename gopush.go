package gopush

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
}

type Interpreter struct {
	Stacks  map[string]*Stack
	Options Options
}

var DefaultOptions = Options{}

func NewInterpreter(options Options) *Interpreter {
	interpreter := &Interpreter{
		Stacks:  make(map[string]*Stack),
		Options: options,
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
