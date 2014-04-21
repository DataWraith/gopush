package gopush

import (
	"math"
	"reflect"
)

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
		p = ignoreWhiteSpace(p, true)
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
			c.Length += 1 + sublist.Length
		} else {
			c.List = append(c.List, Code{Length: 1, Literal: t})
			c.Length++
		}
	}

	return c, nil
}

// Container returns the "container" of the given Code c2 in c. That is, it
// returns the smallest sublist of c which contains c2, or the empty list if
// none of the sublists in c contain c2.
func (c Code) Container(c2 Code) (container Code) {
	container.Length = math.MaxInt32

	if c.Literal != "" {
		return
	}

	for _, sl := range c.List {
		if reflect.DeepEqual(sl, c2) && c.Length < container.Length {
			container = c
		} else {
			candidate := sl.Container(c2)
			if candidate.Length < container.Length {
				container = candidate
			}
		}
	}

	if container.Length == math.MaxInt32 {
		return Code{}
	}

	return container
}

// Contains returns whether the Code c is equal to c2 or contains it in any
// sublist
func (c Code) Contains(c2 Code) bool {
	if reflect.DeepEqual(c, c2) {
		return true
	}

	for _, sl := range c.List {
		if sl.Contains(c2) {
			return true
		}
	}

	return false
}

// UniqueItems returns a map with the count of all unique items in the Code list
func (c Code) UniqueItems() map[string]int64 {
	result := make(map[string]int64)

	if c.Literal != "" {
		result[c.Literal]++
		return result
	}

	for _, sl := range c.List {
		for k, v := range sl.UniqueItems() {
			result[k] += v
		}
	}

	return result
}
