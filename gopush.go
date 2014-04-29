package gopush

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/cryptix/goremutake"
)

// Interpreter is a Push interpreter.
type Interpreter struct {
	Stacks  map[string]*Stack
	Options Options
	Rand    *rand.Rand

	Definitions        map[string]Code
	listOfDefinitions  []string
	listOfInstructions []string

	numEvalPush       int
	quoteNextName     bool
	numNamesGenerated uint
}

// NewInterpreter returns a new Push Interpreter, configured with the provided Options.
func NewInterpreter(options Options) *Interpreter {

	if options.RandomSeed == 0 {
		options.RandomSeed = rand.Int63()
	}

	interpreter := &Interpreter{
		Stacks:             make(map[string]*Stack),
		Options:            options,
		Rand:               rand.New(rand.NewSource(options.RandomSeed)),
		Definitions:        make(map[string]Code),
		listOfDefinitions:  make([]string, 0),
		listOfInstructions: make([]string, 0),
		numEvalPush:        0,
		quoteNextName:      false,
		numNamesGenerated:  0,
	}

	// Setup stacks
	interpreter.RegisterStack("exec", newExecStack(interpreter))
	interpreter.RegisterStack("name", newNameStack(interpreter))
	interpreter.listOfInstructions = append(interpreter.listOfInstructions, "NAME-ERC")

	if _, ok := options.AllowedTypes["boolean"]; ok {
		interpreter.RegisterStack("boolean", newBooleanStack(interpreter))
	}

	if _, ok := options.AllowedTypes["code"]; ok {
		interpreter.RegisterStack("code", newCodeStack(interpreter))
	}

	if _, ok := options.AllowedTypes["float"]; ok {
		interpreter.RegisterStack("float", newFloatStack(interpreter))
		interpreter.listOfInstructions = append(interpreter.listOfInstructions, "FLOAT-ERC")
	}

	if _, ok := options.AllowedTypes["integer"]; ok {
		interpreter.RegisterStack("integer", newIntStack(interpreter))
		interpreter.listOfInstructions = append(interpreter.listOfInstructions, "INTEGER-ERC")
	}

	return interpreter
}

// RegisterStack registers the given stack under the given name. This
// automatically prunes instructions that are not in the set of allowed
// instructions and also makes the instructions of the stack available for
// CODE.RAND to generate. It will NOT overwrite already existing stacks.
func (i *Interpreter) RegisterStack(name string, s *Stack) {
	if _, ok := i.Stacks[name]; ok {
		return
	}

	i.Stacks[name] = s

	// Prune disallowed instructions
	for fn := range s.Functions {
		if _, ok := i.Options.AllowedInstructions[name+"."+fn]; !ok {
			delete(s.Functions, fn)
		}
	}

	// Add the Stack's functions to the list of functions
	for fn := range s.Functions {
		i.listOfInstructions = append(i.listOfInstructions, strings.ToUpper(name+"."+fn))
	}

	// Sort the instructions (otherwise runs aren't repeatable)
	sort.Strings(i.listOfInstructions)
}

func (i *Interpreter) randomInstruction() Code {
	var instr string

	n := i.Rand.Intn(len(i.listOfInstructions) + len(i.listOfDefinitions))

	if n < len(i.listOfInstructions) {
		instr = i.listOfInstructions[n]
	} else {
		instr = i.listOfDefinitions[n-len(i.listOfInstructions)]
	}

	switch instr {
	case "INTEGER-ERC":
		// Generate ephemeral random constant integer
		high := i.Options.MaxRandomInteger
		low := i.Options.MinRandomInteger
		instr = fmt.Sprint(i.Rand.Int63n(high+1-low) + low)

	case "FLOAT-ERC":
		// Generate ephemeral random constant float
		high := i.Options.MaxRandomFloat
		low := i.Options.MinRandomFloat
		instr = fmt.Sprint(i.Rand.Float64()*(high-low) + low)
		if !strings.Contains(instr, ".") {
			instr += ".0"
		}

	case "NAME-ERC":
		// Generate ephemeral random constant NAME
		if i.Rand.Float64() < i.Options.NewERCNameProbabilty || i.numNamesGenerated == 0 {
			// Generate a new random NAME
			instr = goremutake.Encode(i.numNamesGenerated)
			i.numNamesGenerated++
		} else {
			// Use a random, already generated NAME
			instr = goremutake.Encode(uint(i.Rand.Intn(int(i.numNamesGenerated))))
		}
	}

	return Code{Length: 1, Literal: instr}
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

func (i *Interpreter) define(name string, code Code) {
	if _, ok := i.Definitions[name]; !ok {
		i.listOfDefinitions = append(i.listOfDefinitions, name)
	}

	i.Definitions[name] = code
}

func (i *Interpreter) printInterpreterState() {
	fmt.Println("Step", i.numEvalPush)
	for k, v := range i.Stacks {
		fmt.Printf("%s:\n", k)
		for i := len(v.Stack) - 1; i >= 0; i-- {
			fmt.Printf("- %v\n", v.Stack[i])
		}
	}
	fmt.Println()
	fmt.Println()
}

func (i *Interpreter) runCode(program Code) (err error) {

	// Recover from a panic that could occur while executing an instruction.
	// Because it is more convenient for functions to not return an error,
	// the functions that want to return an error panic instead.
	defer func() {
		if perr := recover(); perr != nil {
			err = perr.(error)
		}
	}()

	i.Stacks["exec"].Push(program)

	for i.Stacks["exec"].Len() > 0 && i.numEvalPush < i.Options.EvalPushLimit {

		if i.Options.Tracing {
			i.printInterpreterState()
		}

		item := i.Stacks["exec"].Pop().(Code)
		i.numEvalPush++

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
			if !i.stackOK("integer", 0) {
				return fmt.Errorf("found integer literal %v, but the integer stack is disabled", intlit)
			}
			i.Stacks["integer"].Push(intlit)
			continue
		}

		if floatlit, err := strconv.ParseFloat(item.Literal, 64); err == nil {
			if !i.stackOK("float", 0) {
				return fmt.Errorf("found float literal %v, but the float stack is disabled", floatlit)
			}
			i.Stacks["float"].Push(floatlit)
			continue
		}

		if boollit, err := strconv.ParseBool(item.Literal); err == nil {
			if !i.stackOK("boolean", 0) {
				return fmt.Errorf("found boolean literal %v, but the boolean stack is disabled", boollit)
			}
			i.Stacks["boolean"].Push(boollit)
			continue
		}

		// Try to parse the item on top of the exec stack as instruction
		if strings.Contains(item.Literal, ".") {
			stack := strings.ToLower(item.Literal[:strings.Index(item.Literal, ".")])
			operation := strings.ToLower(item.Literal[strings.Index(item.Literal, ".")+1:])

			s, ok := i.Stacks[stack]
			if !ok {
				return fmt.Errorf("unknown or disabled stack: %v", stack)
			}

			f, ok := s.Functions[operation]
			if !ok {
				return fmt.Errorf("unknown or disabled instruction %v.%v", stack, operation)
			}

			f()
			continue
		}

		// If the item is not an instruction, it must be a name, either
		// bound or unbound. If the quoteNextName flag is false, we can
		// check if the name is already bound.
		if !i.quoteNextName {
			if d, ok := i.Definitions[strings.ToLower(item.Literal)]; ok {
				// Name is already bound, push its value onto the exec stack
				i.Stacks["exec"].Push(d)
				continue
			}
		}

		// The name is not bound yet, so push it onto the name stack
		i.Stacks["name"].Push(strings.ToLower(item.Literal))
		i.quoteNextName = false
	}

	if i.numEvalPush >= i.Options.EvalPushLimit {
		return errors.New("the EvalPushLimit was exceeded")
	}

	return nil
}

// RunCode runs the given program (given as Code type) until the EvalPushLimit
// is reached
func (i *Interpreter) RunCode(c Code) error {
	if i.Options.TopLevelPushCode {
		if s, ok := i.Stacks["code"]; ok {
			s.Push(c)
		}
	}

	err := i.runCode(c)

	if i.Options.TopLevelPopCode {
		if s, ok := i.Stacks["code"]; ok {
			s.Pop()
		}
	}

	if i.Options.Tracing {
		i.printInterpreterState()
	}

	return err
}

// Run runs the given program written in the Push programming language until the
// EvalPushLimit is reached
func (i *Interpreter) Run(program string) error {
	c, err := ParseCode(program)
	if err != nil {
		return err
	}

	err = i.RunCode(c)

	return err
}
