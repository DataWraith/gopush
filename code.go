package gopush

// Code is the internal list representation of a (partial) Push program.
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

// ParseCode takes the provided Push program and parses it into the internal
// list representation (type Code).
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
