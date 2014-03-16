package gopush

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Options struct {
	// When TRUE (which is the default), code passed to the top level of
	// the interpreter will be pushed onto the CODE stack prior to
	// execution.
	TopLevelPushCode bool

	// When TRUE, the CODE stack will be popped at the end of top level
	// calls to the interpreter. The default is FALSE.
	TopLevelPopCode bool

	// The maximum number of points that will be executed in a single
	// top-level call to the interpreter.
	EvalPushLimit int

	// The probability that the selection of the ephemeral random NAME
	// constant for inclusion in randomly generated code will produce a new
	// name (rather than a name that was previously generated).
	NewERCNameProbabilty float64

	// The maximum number of points that can occur in any program on the
	// CODE stack. Instructions that would violate this limit act as NOOPs.
	MaxPointsInProgram int

	// The maximum number of points in an expression produced by the
	// CODE.RAND instruction.
	MaxPointsInRandomExpression int

	// The maximum FLOAT that will be produced as an ephemeral random FLOAT
	// constant or from a call to FLOAT.RAND.
	MaxRandomFloat float64

	// The minimum FLOAT that will be produced as an ephemeral random FLOAT
	// constant or from a call to FLOAT.RAND.
	MinRandomFloat float64

	// The maximum INTEGER that will be produced as an ephemeral random
	// INTEGER constant or from a call to INTEGER.RAND.
	MaxRandomInteger int64

	// The minimum INTEGER that will be produced as an ephemeral random
	// INTEGER constant or from a call to INTEGER.RAND.
	MinRandomInteger int64

	// A seed for the random number generator.
	RandomSeed int64
}

type Definition struct {
	Stack string
	Value interface{}
}

type Interpreter struct {
	Stacks      map[string]*Stack
	Options     Options
	Rand        *rand.Rand
	Definitions map[string]Definition
}

var DefaultOptions = Options{
	TopLevelPushCode:            true,
	TopLevelPopCode:             false,
	EvalPushLimit:               1000,
	NewERCNameProbabilty:        0.001,
	MaxPointsInProgram:          100,
	MaxPointsInRandomExpression: 25,
	MaxRandomFloat:              1.0,
	MinRandomFloat:              -1.0,
	MaxRandomInteger:            10,
	MinRandomInteger:            -10,
	RandomSeed:                  rand.Int63(),
}

func NewInterpreter(options Options) *Interpreter {
	interpreter := &Interpreter{
		Stacks:      make(map[string]*Stack),
		Options:     options,
		Rand:        rand.New(rand.NewSource(options.RandomSeed)),
		Definitions: make(map[string]Definition),
	}

	interpreter.Stacks["integer"] = NewIntStack(interpreter)
	interpreter.Stacks["float"] = NewFloatStack(interpreter)
	interpreter.Stacks["exec"] = new(Stack)
	interpreter.Stacks["code"] = NewCodeStack(interpreter)
	interpreter.Stacks["name"] = new(Stack)
	interpreter.Stacks["boolean"] = NewBooleanStack(interpreter)

	return interpreter
}

func (i *Interpreter) stackOK(name string, mindepth int64) bool {
	s, ok := i.Stacks[name]
	if !ok {
		return false
	}

	if s.Len() < mindepth {
		return false
	}

	return true
}

func (i *Interpreter) Run(program string) error {
	c, err := ParseCode(program)
	if err != nil {
		return err
	}

	i.Stacks["exec"].Push(c)

	numEvalPush := 0

	for i.Stacks["exec"].Len() > 0 && numEvalPush < i.Options.EvalPushLimit {
		item := i.Stacks["exec"].Pop().(Code)
		numEvalPush++

		// If the item on top of the exec stack is a list, push it in
		// reverse order
		if item.Literal == "" {
			for j := len(item.List) - 1; j >= 0; j-- {
				i.Stacks["exec"].Push(item.List[j])
			}
			continue
		}

		// Try to parse the item on top of the exec stack as a literal
		if intlit, err := strconv.ParseInt(item.Literal, 10, 64); err == nil {
			i.Stacks["integer"].Push(intlit)
			continue
		}

		if floatlit, err := strconv.ParseFloat(item.Literal, 64); err == nil {
			i.Stacks["float"].Push(floatlit)
			continue
		}

		if boollit, err := strconv.ParseBool(item.Literal); err == nil {
			i.Stacks["boolean"].Push(boollit)
			continue
		}

		// Try to parse the item on top of the exec stack as instruction
		if strings.Contains(item.Literal, ".") {
			stack := strings.ToLower(item.Literal[:strings.Index(item.Literal, ".")])
			operation := strings.ToLower(item.Literal[strings.Index(item.Literal, ".")+1:])

			s, ok := i.Stacks[stack]
			if !ok {
				return errors.New(fmt.Sprintf("unkown or disabled stack: %v", stack))
			}

			f, ok := s.Functions[operation]
			if !ok {
				return errors.New(fmt.Sprintf("unknown or disabled instruction: %v.%v", stack, operation))
			}

			f()
			continue
		}

		// If the item is not an instruction, it must be a name, either
		// bound or unbound. First we check for bound names.
		d, ok := i.Definitions[strings.ToLower(item.Literal)]
		if ok {
			// Name is already bound, push the bound value onto the appropriate stack
			i.Stacks[d.Stack].Push(d.Value)
			continue
		}

		// The item is not bound yet, so push it onto the name stack
		i.Stacks["name"].Push(strings.ToLower(item.Literal))
	}

	return nil
}
