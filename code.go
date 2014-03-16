package gopush

import (
	"errors"
	"unicode"
)

type Code struct {
	Length  int
	Literal string
	List    []Code
}

func (c Code) String() string {
	if c.Literal != "" {
		return c.Literal
	}

	s := "( "
	for _, v := range c.List {
		s += v.String() + " "
	}
	return s + ")"
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
	return "", "", errors.New("unbalanced parentheses")
}

func ParseCode(program string) (c Code, err error) {
	t := ""
	p := program

	for len(p) > 0 {
		p = ignoreWhiteSpace(p)
		t, p = getToken(p)

		if t == "" {
			break
		}

		if t == "(" {
			t, p, err = getToParen(p)
			if err != nil {
				return Code{}, err
			}

			sublist, err := ParseCode(t)
			if err != nil {
				return Code{}, err
			}

			c.List = append(c.List, sublist)
			c.Length += sublist.Length
		} else {
			c.List = append(c.List, Code{Length: 1, Literal: t})
			c.Length++
		}
	}

	return c, nil
}
