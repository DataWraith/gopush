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
