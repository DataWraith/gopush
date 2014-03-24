package gopush

import (
	"errors"
	"unicode"
)

func ignoreWhiteSpace(program string) string {
	for i, r := range program {
		if !unicode.IsSpace(r) {
			return program[i:]
		}
	}
	return ""
}

func ignoreComments(s string) string {
	for {
		s = ignoreWhiteSpace(s)
		if s[0] == '#' {
			for i, r := range s {
				if r == '\n' {
					s = s[i+1:]
					break
				}
			}
		} else {
			return s
		}
	}
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

func getParameterSettingPair(s string) (parameter, setting, remainder string) {
	s = ignoreComments(s)
	if s == "" {
		return "", "", ""
	}

	parameter, s = getToken(s)
	s = ignoreWhiteSpace(s)
	setting, s = getToken(s)

	return parameter, setting, s
}
