package gopush

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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

	// When TRUE the interpreter will print out the stacks after every
	// executed instruction
	Tracing bool

	// A seed for the random number generator.
	RandomSeed int64
}

// DefaultOptions hold the default options for the Push Interpreter
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
	Tracing:                     false,
	RandomSeed:                  rand.Int63(),
}

func parseOptions(s string) (Options, error) {
	o := DefaultOptions

	var parameter, setting string

	for len(s) > 0 {

		parameter, setting, s = getParameterSettingPair(s)

		if parameter == "" {
			break
		}

		switch strings.ToLower(parameter) {
		case "type":
		case "instruction":
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
				return Options{}, fmt.Errorf("MAX-POINTS-IN-RANDOM-EXPRESSION must be at least 1, got %v", i)
			}

			o.MaxPointsInRandomExpression = int(i)

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
