package gopush

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// Options holds the configuration options for a Push Interpreter
type Options struct {
	// When TRUE (which is the default), code passed to the top level of the
	// interpreter will be pushed onto the CODE stack prior to execution.
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
	MaxPointsInRandomExpression int64

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

	// When TRUE the interpreter will print out the stacks after every
	// executed instruction
	Tracing bool

	// A seed for the random number generator.
	RandomSeed int64

	// AllowedTypes lists the types (stacks) that are allowed
	AllowedTypes map[string]struct{}

	// AllowedInstructions lists the instructions that are allowed
	AllowedInstructions map[string]struct{}
}

// DefaultConfigFile holds the textual representation of the default
// configuration
var DefaultConfigFile = `
## PARAMETER SETTINGS
TOP-LEVEL-PUSH-CODE TRUE
TOP-LEVEL-POP-CODE FALSE

EVALPUSH-LIMIT 1000

NEW-ERC-NAME-PROBABILITY 0.001

MAX-POINTS-IN-PROGRAM 100
MAX-POINTS-IN-RANDOM-EXPRESSIONS 25

MAX-RANDOM-FLOAT 1.0
MIN-RANDOM-FLOAT -1.0

MAX-RANDOM-INTEGER 10
MIN-RANDOM-INTEGER -10

TRACING FALSE


## TYPES
type BOOLEAN
type CODE
type FLOAT
type INTEGER


## INSTRUCTIONS
instruction BOOLEAN.=
instruction BOOLEAN.AND
instruction BOOLEAN.DEFINE
instruction BOOLEAN.DUP
instruction BOOLEAN.FLUSH
instruction BOOLEAN.FROMFLOAT
instruction BOOLEAN.FROMINTEGER
instruction BOOLEAN.NOT
instruction BOOLEAN.OR
instruction BOOLEAN.POP
instruction BOOLEAN.RAND
instruction BOOLEAN.ROT
instruction BOOLEAN.SHOVE
instruction BOOLEAN.STACKDEPTH
instruction BOOLEAN.SWAP
instruction BOOLEAN.YANK
instruction BOOLEAN.YANKDUP

instruction CODE.=
instruction CODE.APPEND
instruction CODE.ATOM
instruction CODE.CAR
instruction CODE.CDR
instruction CODE.CONS
instruction CODE.CONTAINER
instruction CODE.CONTAINS
instruction CODE.DEFINE
instruction CODE.DEFINITION
instruction CODE.DISCREPANCY
instruction CODE.DO
instruction CODE.DO*
instruction CODE.DO*COUNT
instruction CODE.DO*RANGE
instruction CODE.DO*TIMES
instruction CODE.DUP
instruction CODE.EXTRACT
instruction CODE.FLUSH
instruction CODE.FROMBOOLEAN
instruction CODE.FROMFLOAT
instruction CODE.FROMINTEGER
instruction CODE.FROMNAME
instruction CODE.IF
instruction CODE.INSERT
instruction CODE.INSTRUCTIONS
instruction CODE.LENGTH
instruction CODE.LIST
instruction CODE.MEMBER
instruction CODE.NOOP
instruction CODE.NTH
instruction CODE.NTHCDR
instruction CODE.NULL
instruction CODE.POP
instruction CODE.POSITION
instruction CODE.QUOTE
instruction CODE.RAND
instruction CODE.ROT
instruction CODE.SHOVE
instruction CODE.SIZE
instruction CODE.STACKDEPTH
instruction CODE.SUBST
instruction CODE.SWAP
instruction CODE.YANK
instruction CODE.YANKDUP

instruction EXEC.=
instruction EXEC.DEFINE
instruction EXEC.DO*COUNT
instruction EXEC.DO*RANGE
instruction EXEC.DO*TIMES
instruction EXEC.DUP
instruction EXEC.FLUSH
instruction EXEC.IF
instruction EXEC.K
instruction EXEC.POP
instruction EXEC.ROT
instruction EXEC.S
instruction EXEC.SHOVE
instruction EXEC.STACKDEPTH
instruction EXEC.SWAP
instruction EXEC.Y
instruction EXEC.YANK
instruction EXEC.YANKDUP

instruction FLOAT.%
instruction FLOAT.*
instruction FLOAT.+
instruction FLOAT.-
instruction FLOAT./
instruction FLOAT.<
instruction FLOAT.=
instruction FLOAT.>
instruction FLOAT.COS
instruction FLOAT.DEFINE
instruction FLOAT.DUP
instruction FLOAT.FLUSH
instruction FLOAT.FROMBOOLEAN
instruction FLOAT.FROMINTEGER
instruction FLOAT.MAX
instruction FLOAT.MIN
instruction FLOAT.POP
instruction FLOAT.RAND
instruction FLOAT.ROT
instruction FLOAT.SHOVE
instruction FLOAT.SIN
instruction FLOAT.STACKDEPTH
instruction FLOAT.SWAP
instruction FLOAT.TAN
instruction FLOAT.YANK
instruction FLOAT.YANKDUP

instruction INTEGER.%
instruction INTEGER.*
instruction INTEGER.+
instruction INTEGER.-
instruction INTEGER./
instruction INTEGER.<
instruction INTEGER.=
instruction INTEGER.>
instruction INTEGER.DEFINE
instruction INTEGER.DUP
instruction INTEGER.FLUSH
instruction INTEGER.FROMBOOLEAN
instruction INTEGER.FROMFLOAT
instruction INTEGER.MAX
instruction INTEGER.MIN
instruction INTEGER.POP
instruction INTEGER.RAND
instruction INTEGER.ROT
instruction INTEGER.SHOVE
instruction INTEGER.STACKDEPTH
instruction INTEGER.SWAP
instruction INTEGER.YANK
instruction INTEGER.YANKDUP

instruction NAME.=
instruction NAME.DUP
instruction NAME.FLUSH
instruction NAME.POP
instruction NAME.QUOTE
instruction NAME.RAND
instruction NAME.RANDBOUNDNAME
instruction NAME.ROT
instruction NAME.SHOVE
instruction NAME.STACKDEPTH
instruction NAME.SWAP
instruction NAME.YANK
instruction NAME.YANKDUP
`

// DefaultOptions contains the default configuration for a Push Interpreter.
var DefaultOptions, _ = parseOptions(DefaultConfigFile)

func parseOptions(s string) (Options, error) {
	o := Options{
		AllowedInstructions:         make(map[string]struct{}),
		AllowedTypes:                make(map[string]struct{}),
		EvalPushLimit:               1000,
		MaxPointsInProgram:          100,
		MaxPointsInRandomExpression: 25,
		MaxRandomFloat:              1.0,
		MaxRandomInteger:            10,
		MinRandomFloat:              -1.0,
		MinRandomInteger:            -10,
		NewERCNameProbabilty:        0.001,
		RandomSeed:                  0,
		TopLevelPopCode:             false,
		TopLevelPushCode:            true,
		Tracing:                     false,
	}

	var parameter, setting string

	for len(s) > 0 {

		parameter, setting, s = getParameterSettingPair(s)

		if parameter == "" {
			break
		}

		if setting == "" {
			return Options{}, fmt.Errorf("expected setting to follow %q", parameter)
		}

		switch strings.ToLower(parameter) {
		case "type":
			t := strings.ToLower(setting)
			switch t {
			case "boolean":
				fallthrough
			case "code":
				fallthrough
			case "float":
				fallthrough
			case "integer":
				o.AllowedTypes[t] = struct{}{}

			// NAME and EXEC stacks always exist, so they are a
			// no-op with the type parameter
			case "name":
			case "exec":

			default:
				return Options{}, fmt.Errorf("unknown type: %q", setting)
			}

		case "instruction":
			instr := strings.ToLower(setting)
			o.AllowedInstructions[instr] = struct{}{}

		case "min-random-integer":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}
			o.MinRandomInteger = i

		case "max-random-integer":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}
			o.MinRandomInteger = i

		case "min-random-float":
			f, err := strconv.ParseFloat(setting, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as float", setting)
			}
			o.MinRandomFloat = f

		case "max-random-float":
			f, err := strconv.ParseFloat(setting, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as float", setting)
			}
			o.MaxRandomFloat = f

		case "max-points-in-random-expressions":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}

			if i < 1 {
				return Options{}, fmt.Errorf("MAX-POINTS-IN-RANDOM-EXPRESSIONS must be at least 1, got %v", i)
			}

			o.MaxPointsInRandomExpression = i

		case "max-points-in-program":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}

			if i < 1 {
				return Options{}, fmt.Errorf("MAX-POINTS-IN-PROGRAM must be at least 1, got %v", i)
			}

			o.MaxPointsInProgram = int(i)

		case "evalpush-limit":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}

			if i < 1 {
				return Options{}, fmt.Errorf("EVALPUSH-LIMIT must be at least 1, got %v", i)
			}

			o.EvalPushLimit = int(i)

		case "new-erc-name-probability":
			f, err := strconv.ParseFloat(setting, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as float", setting)
			}

			if f < 0 || f > 1 {
				return Options{}, fmt.Errorf("NEW-ERC-NAME-PROBABILITY must be between 0 and 1 inclusive, got %v", f)
			}

			o.NewERCNameProbabilty = f

		case "random-seed":
			i, err := strconv.ParseInt(setting, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as integer", setting)
			}

			o.RandomSeed = i

		case "top-level-push-code":
			b, err := strconv.ParseBool(setting)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as boolean", setting)
			}

			o.TopLevelPushCode = b

		case "top-level-pop-code":
			b, err := strconv.ParseBool(setting)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as boolean", setting)
			}

			o.TopLevelPopCode = b

		case "tracing":
			b, err := strconv.ParseBool(setting)
			if err != nil {
				return Options{}, fmt.Errorf("could not parse %q as boolean", setting)
			}

			o.Tracing = b
		default:
			return Options{}, fmt.Errorf("unknown parameter %q", parameter)
		}
	}

	return o, nil
}

// ReadOptions reads the options from the given reader
func ReadOptions(r io.Reader) (Options, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return Options{}, err
	}

	s := string(b)

	o, err := parseOptions(s)
	if err != nil {
		return Options{}, err
	}

	return o, nil
}
