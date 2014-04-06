package gopush

import (
	"fmt"
	"reflect"
)

func newExecStack(interpreter *Interpreter) *Stack {
	s := &Stack{
		Functions: make(map[string]Instruction),
	}

	s.Functions["="] = func() {
		if !interpreter.stackOK("exec", 2) || !interpreter.stackOK("boolean", 0) {
			return
		}

		e1 := interpreter.Stacks["exec"].Pop().(Code)
		e2 := interpreter.Stacks["exec"].Pop().(Code)
		same := reflect.DeepEqual(e1, e2)
		interpreter.Stacks["boolean"].Push(same)
	}

	s.Functions["define"] = func() {
		if !interpreter.stackOK("exec", 1) || !interpreter.stackOK("name", 1) {
			return
		}

		n := interpreter.Stacks["name"].Pop().(string)
		c := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.define(n, c)
	}

	s.Functions["do*count"] = func() {
		if !interpreter.stackOK("exec", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		if _, ok := interpreter.Stacks["exec"].Functions["do*range"]; !ok {
			return
		}

		count := interpreter.Stacks["integer"].Pop().(int64)
		code := interpreter.Stacks["exec"].Pop().(Code)

		if count <= 0 {
			return
		}

		toPush := Code{
			Length: 3 + code.Length,
			List: []Code{
				Code{Length: 1, Literal: "0"},
				Code{Length: 1, Literal: fmt.Sprint(count - 1)},
				Code{Length: 1, Literal: "EXEC.DO*RANGE"},
				code,
			}}

		interpreter.Stacks["exec"].Push(toPush)
	}

	s.Functions["do*range"] = func() {
		if !interpreter.stackOK("exec", 1) || !interpreter.stackOK("integer", 2) {
			return
		}

		c := interpreter.Stacks["exec"].Pop().(Code)
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

			interpreter.Stacks["exec"].Push(Code{
				Length: 3 + c.Length,
				List: []Code{
					Code{Length: 1, Literal: fmt.Sprint(cur)},
					Code{Length: 1, Literal: fmt.Sprint(dst)},
					Code{Length: 1, Literal: "EXEC.DO*RANGE"},
					c,
				},
			})

			interpreter.Stacks["exec"].Push(c)
		}
	}

	s.Functions["do*times"] = func() {
		if !interpreter.stackOK("exec", 1) || !interpreter.stackOK("integer", 1) {
			return
		}

		if _, ok := interpreter.Stacks["exec"].Functions["do*range"]; !ok {
			return
		}

		count := interpreter.Stacks["integer"].Pop().(int64)
		code := interpreter.Stacks["exec"].Pop().(Code)

		if count <= 0 {
			return
		}

		loopBody := Code{
			Length: 1 + code.Length,
			List: []Code{
				Code{Length: 1, Literal: "INTEGER.POP"},
				code,
			},
		}

		toPush := Code{
			Length: 3 + code.Length,
			List: []Code{
				Code{Length: 1, Literal: "0"},
				Code{Length: 1, Literal: fmt.Sprint(count - 1)},
				Code{Length: 1, Literal: "EXEC.DO*RANGE"},
				loopBody,
			}}

		interpreter.Stacks["exec"].Push(toPush)
	}

	s.Functions["dup"] = func() {
		interpreter.Stacks["exec"].Dup()
	}

	s.Functions["flush"] = func() {
		interpreter.Stacks["exec"].Flush()
	}

	s.Functions["if"] = func() {
		if !interpreter.stackOK("exec", 2) || !interpreter.stackOK("boolean", 1) {
			return
		}

		b := interpreter.Stacks["boolean"].Pop().(bool)
		c1 := interpreter.Stacks["exec"].Pop().(Code)
		c2 := interpreter.Stacks["exec"].Pop().(Code)

		if b {
			interpreter.Stacks["exec"].Push(c1)
		} else {
			interpreter.Stacks["exec"].Push(c2)
		}
	}

	s.Functions["k"] = func() {
		if !interpreter.stackOK("exec", 2) {
			return
		}

		i1 := interpreter.Stacks["exec"].Pop().(Code)
		_ = interpreter.Stacks["exec"].Pop()
		interpreter.Stacks["exec"].Push(i1)
	}

	s.Functions["pop"] = func() {
		interpreter.Stacks["exec"].Pop()
	}

	s.Functions["rot"] = func() {
		interpreter.Stacks["exec"].Rot()
	}

	s.Functions["s"] = func() {
		if !interpreter.stackOK("exec", 3) {
			return
		}

		a := interpreter.Stacks["exec"].Pop().(Code)
		b := interpreter.Stacks["exec"].Pop().(Code)
		c := interpreter.Stacks["exec"].Pop().(Code)

		l := Code{
			Length: b.Length + c.Length,
			List: []Code{
				b,
				c,
			},
		}

		interpreter.Stacks["exec"].Push(l)
		interpreter.Stacks["exec"].Push(c)
		interpreter.Stacks["exec"].Push(a)
	}

	s.Functions["shove"] = func() {
		if !interpreter.stackOK("integer", 1) || !interpreter.stackOK("exec", 1) {
			return
		}

		i := interpreter.Stacks["integer"].Pop().(int64)
		c := interpreter.Stacks["exec"].Peek().(Code)

		interpreter.Stacks["exec"].Shove(c, i)
		interpreter.Stacks["exec"].Pop()
	}

	s.Functions["stackdepth"] = func() {
		if !interpreter.stackOK("integer", 0) {
			return
		}

		interpreter.Stacks["integer"].Push(interpreter.Stacks["exec"].Len())
	}

	s.Functions["y"] = func() {
		if !interpreter.stackOK("exec", 1) {
			return
		}

		e := interpreter.Stacks["exec"].Pop().(Code)
		interpreter.Stacks["exec"].Push(Code{Length: 1 + e.Length, List: []Code{Code{Length: 1, Literal: "EXEC.Y"}, e}})
		interpreter.Stacks["exec"].Push(e)
	}

	return s
}
