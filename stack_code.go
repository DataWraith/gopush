package gopush

import (
	"fmt"
	"reflect"
	"strings"
)

func newCodeStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]func()),
	}

	s.Functions["="] = func() {
		if !interpreter.StackOK("code", 2) || !interpreter.StackOK("boolean", 0) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if reflect.DeepEqual(c1, c2) {
			interpreter.Stacks["boolean"].Push(true)
		} else {
			interpreter.Stacks["boolean"].Push(false)
		}
	}

	s.Functions["append"] = func() {
		if !interpreter.StackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if c1.Literal != "" {
			c1 = Code{Length: c1.Length, List: []Code{c1}}
		}

		if c2.Literal != "" {
			c2 = Code{Length: c2.Length, List: []Code{c2}}
		}

		combined := Code{Length: c1.Length + c2.Length, List: append(c2.List, c1.List...)}

		if combined.Length <= interpreter.Options.MaxPointsInProgram {
			interpreter.Stacks["code"].Push(combined)
		}
	}

	s.Functions["atom"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("boolean", 0) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if c.Literal != "" {
			interpreter.Stacks["boolean"].Push(true)
		} else {
			interpreter.Stacks["boolean"].Push(false)
		}
	}

	s.Functions["car"] = func() {
		if !interpreter.StackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if len(c.List) == 0 {
			return
		}

		interpreter.Stacks["code"].Push(c.List[0])
	}

	s.Functions["cdr"] = func() {
		if !interpreter.StackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		if len(c.List) == 0 {
			interpreter.Stacks["code"].Push(Code{})
		} else {
			cdr := Code{
				Length: c.Length - c.List[0].Length,
				List:   c.List[1:],
			}
			interpreter.Stacks["code"].Push(cdr)
		}
	}

	s.Functions["cons"] = func() {
		if !interpreter.StackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if c1.Literal != "" {
			c1 = Code{Length: 1, List: []Code{c1}}
		}

		if c2.Literal != "" {
			c2 = Code{Length: 1, List: []Code{c2}}
		}

		c := Code{
			Length: c1.Length + c2.Length,
			List:   append(c2.List, c1.List...),
		}

		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["container"] = func() {
		if !interpreter.StackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		c := c1.Container(c2)
		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["contains"] = func() {
		if !interpreter.StackOK("code", 2) || !interpreter.StackOK("boolean", 0) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		interpreter.Stacks["boolean"].Push(c2.Contains(c1))
	}

	s.Functions["define"] = func() {
		if !interpreter.StackOK("name", 1) || !interpreter.StackOK("code", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		c := interpreter.Stacks["code"].Pop().(Code)

		interpreter.define(n, c)
	}

	s.Functions["definition"] = func() {
		if !interpreter.StackOK("name", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)

		if c, ok := interpreter.Definitions[n]; ok {
			interpreter.Stacks["code"].Push(c)
		}
	}

	s.Functions["discrepancy"] = func() {
		if !interpreter.StackOK("code", 2) || !interpreter.StackOK("integer", 0) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		u1 := c1.UniqueItems()
		u2 := c2.UniqueItems()

		keys := make(map[string]struct{})

		for k := range u1 {
			keys[k] = struct{}{}
		}

		for k := range u2 {
			keys[k] = struct{}{}
		}

		discrepancy := int64(0)
		for k := range keys {
			if u1[k] > u2[k] {
				discrepancy += u1[k] - u2[k]
			} else {
				discrepancy += u2[k] - u1[k]
			}
		}

		interpreter.Stacks["integer"].Push(discrepancy)
	}

	s.Functions["do"] = func() {
		if !interpreter.StackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)

		err := interpreter.runCode(c)
		if err != nil {
			panic(err)
		}

		interpreter.Stacks["code"].Pop()
	}

	s.Functions["do*"] = func() {
		if !interpreter.StackOK("code", 1) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)
		interpreter.Stacks["code"].Pop()

		err := interpreter.runCode(c)
		if err != nil {
			panic(err)
		}
	}

	s.Functions["do*count"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["code"].Pop().(Code)

		if i <= 0 {
			return
		}

		toPush := Code{
			Length: 4 + c.Length,
			List: []Code{
				Code{Length: 1, Literal: "0"},
				Code{Length: 1, Literal: fmt.Sprint(i)},
				Code{Length: 1, Literal: "CODE.QUOTE"},
				c,
				Code{Length: 1, Literal: "CODE.DO*RANGE"},
			},
		}

		interpreter.Stacks["code"].Push(toPush)
	}

	s.Functions["do*range"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 2) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)
		dst := interpreter.Stacks["integer"].Pop().(int64)
		cur := interpreter.Stacks["integer"].Pop().(int64)

		if cur == dst {
			interpreter.Stacks["integer"].Push(cur)
			interpreter.Stacks["exec"].Push(c)
		} else {
			interpreter.Stacks["integer"].Push(cur)

			if dst < cur {
				cur--
			} else {
				cur++
			}

			interpreter.Stacks["code"].Push(c)
			interpreter.Stacks["exec"].Push(c)
			interpreter.Stacks["exec"].Push(Code{Length: 1, Literal: "CODE.DO*RANGE"})
			interpreter.Stacks["integer"].Push(cur)
			interpreter.Stacks["integer"].Push(dst)
		}
	}

	s.Functions["do*times"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["code"].Pop().(Code)

		if i <= 0 {
			return
		}

		toPush := Code{
			Length: 4 + c.Length,
			List: []Code{
				Code{Length: 1, Literal: "0"},
				Code{Length: 1, Literal: fmt.Sprint(i)},
				Code{Length: 1, Literal: "CODE.QUOTE"},
				Code{Length: 1 + c.Length, List: []Code{
					Code{Length: 1, Literal: "INTEGER.POP"},
					c,
				}},
				Code{Length: 1, Literal: "CODE.DO*RANGE"},
			},
		}

		interpreter.Stacks["code"].Push(toPush)
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["code"].Dup()
	}

	s.Functions["extract"] = func() {
		// TODO
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["code"].Flush()
	}

	s.Functions["fromboolean"] = func() {
		if !interpreter.StackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		interpreter.Stacks["code"].Push(Code{Length: 1, Literal: fmt.Sprint(b)})
	}

	s.Functions["fromfloat"] = func() {
		if !interpreter.StackOK("float", 1) {
			return
		}

		f := interpreter.Stacks["float"].Pop().(float64)
		l := fmt.Sprint(f)
		if !strings.Contains(l, ".") {
			l += ".0"
		}
		interpreter.Stacks["code"].Push(Code{Length: 1, Literal: l})
	}

	s.Functions["frominteger"] = func() {
		if !interpreter.StackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["code"].Push(Code{Length: 1, Literal: fmt.Sprint(i)})
	}

	s.Functions["fromname"] = func() {
		if !interpreter.StackOK("name", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		interpreter.Stacks["code"].Push(Code{Length: 1, Literal: n})
	}

	s.Functions["if"] = func() {
		if !interpreter.StackOK("code", 2) || !interpreter.StackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		if b {
			interpreter.Stacks["exec"].Push(c2)
		} else {
			interpreter.Stacks["exec"].Push(c1)
		}
	}

	s.Functions["insert"] = func() {
		// TODO
	}

	s.Functions["instructions"] = func() {
		c := Code{List: make([]Code, 0, len(interpreter.listOfInstructions))}

		for _, instr := range interpreter.listOfInstructions {
			if instr == "NAME-ERC" || instr == "FLOAT-ERC" || instr == "INTEGER-ERC" {
				continue
			}

			c.Length++
			c.List = append(c.List, Code{Length: 1, Literal: instr})
		}

		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["length"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 0) {
			return
		}

		c := interpreter.Stacks["code"].Peek().(Code)
		if c.Literal != "" {
			interpreter.Stacks["integer"].Push(int64(1))
		} else {
			interpreter.Stacks["integer"].Push(int64(len(c.List)))
		}
	}

	s.Functions["list"] = func() {
		if !interpreter.StackOK("code", 2) {
			return
		}

		c1 := interpreter.Stacks["code"].Pop().(Code)
		c2 := interpreter.Stacks["code"].Pop().(Code)

		c := Code{
			Length: c1.Length + c2.Length,
			List:   []Code{c1, c2},
		}

		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["member"] = func() {
		// TODO
	}

	s.Functions["noop"] = func() {
		// Does nothing
	}

	s.Functions["nth"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["code"].Pop().(Code)

		if c.Literal == "" && len(c.List) == 0 {
			interpreter.Stacks["code"].Push(c)
			return
		}

		if c.Literal != "" {
			c = Code{Length: c.Length, List: []Code{c}}
		}

		idx := i % int64(len(c.List))
		if idx < 0 {
			idx = -idx
		}

		interpreter.Stacks["code"].Push(c.List[idx])
	}

	s.Functions["nthcdr"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["code"].Pop().(Code)

		if c.Literal == "" && len(c.List) == 0 {
			interpreter.Stacks["code"].Push(c)
			return
		}

		if c.Literal != "" {
			c = Code{Length: c.Length, List: []Code{c}}
		}

		idx := i % int64(len(c.List))
		if idx < 0 {
			idx = -idx
		}

		nthcdr := Code{List: c.List[idx:]}
		for _, sl := range c.List {
			nthcdr.Length += sl.Length
		}

		interpreter.Stacks["code"].Push(nthcdr)
	}

	s.Functions["null"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("boolean", 0) {
			return
		}

		c := interpreter.Stacks["code"].Pop().(Code)
		interpreter.Stacks["boolean"].Push(c.Literal == "" && len(c.List) == 0)
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["code"].Pop()
	}

	s.Functions["position"] = func() {
		// TODO
	}

	s.Functions["quote"] = func() {
		if !interpreter.StackOK("exec", 1) {
			return
		}

		c := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["rand"] = func() {
		if !interpreter.StackOK("integer", 1) {
			return
		}

		maxPoints := interpreter.Stacks["integer"].Pop().(int64)

		if maxPoints == 0 {
			return
		}

		if maxPoints < 0 {
			maxPoints *= -1
		}

		if maxPoints > interpreter.Options.MaxPointsInRandomExpression {
			maxPoints = interpreter.Options.MaxPointsInRandomExpression
		}

		c := interpreter.RandomCode(maxPoints)
		interpreter.Stacks["code"].Push(c)
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["code"].Rot()
	}

	s.Functions["shove"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["code"].Peek().(Code)
		interpreter.Stacks["code"].Shove(c, idx)
		interpreter.Stacks["code"].Pop()
	}

	s.Functions["size"] = func() {
		if !interpreter.StackOK("code", 1) || !interpreter.StackOK("integer", 0) {
			return
		}

		c := interpreter.Stacks["code"].Peek().(Code)
		interpreter.Stacks["integer"].Push(c.Length)
	}

	s.Functions["stackdepth"] = func() {
		if !interpreter.StackOK("integer", 0) {
			return
		}

		interpreter.Stacks["integer"].Push(interpreter.Stacks["code"].Len())
	}

	s.Functions["subst"] = func() {
		// TODO
	}

	s.Functions["swap"] = func() {
		interpreter.Stacks["code"].Swap()
	}

	s.Functions["yank"] = func() {
		if !interpreter.StackOK("integer", 1) || !interpreter.StackOK("code", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["code"].Yank(idx)
	}

	s.Functions["yankdup"] = func() {
		if !interpreter.StackOK("integer", 1) || !interpreter.StackOK("code", 1) {
			return
		}

		idx := interpreter.Stacks["integer"].Pop().(int64)
		interpreter.Stacks["code"].YankDup(idx)
	}

	return s
}
