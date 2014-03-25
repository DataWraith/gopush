package gopush_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/DataWraith/gopush"
)

func TestParsingEmptyConfigurationFile(t *testing.T) {
	opt, err := gopush.ReadOptions(strings.NewReader(""))
	if err != nil {
		t.Errorf("unexpected error while parsing configuration: %v", err)
	}

	defaultConfig := gopush.Options{
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

	if !reflect.DeepEqual(opt, defaultConfig) {
		t.Error("expected \"\" to parse to defaultConfig")
	}
}

var optionParseErrorTests = []struct {
	toParse       string
	expectedError string
}{
	{"type foo", "unknown type: \"foo\""},
	{"type \ninteger", "expected setting to follow \"type\""},
	{"min-random-integer foo", "could not parse \"foo\" as integer"},
	{"max-random-integer foo", "could not parse \"foo\" as integer"},
	{"max-points-in-random-expressions foo", "could not parse \"foo\" as integer"},
	{"max-points-in-program foo", "could not parse \"foo\" as integer"},
	{"evalpush-limit foo", "could not parse \"foo\" as integer"},
	{"random-seed foo", "could not parse \"foo\" as integer"},
	{"min-random-float foo", "could not parse \"foo\" as float"},
	{"max-random-float foo", "could not parse \"foo\" as float"},
	{"new-erc-name-probability foo", "could not parse \"foo\" as float"},
	{"top-level-push-code foo", "could not parse \"foo\" as boolean"},
	{"top-level-pop-code foo", "could not parse \"foo\" as boolean"},
	{"tracing foo", "could not parse \"foo\" as boolean"},
	{"foo bar", "unknown parameter \"foo\""},
}

func TestParseErrors(t *testing.T) {
	for _, pe := range optionParseErrorTests {
		_, err := gopush.ReadOptions(strings.NewReader(pe.toParse))
		if err.Error() != pe.expectedError {
			t.Errorf("unexpected error while parsing configuration: %q, expected %q", err, pe.expectedError)
		}
	}
}
