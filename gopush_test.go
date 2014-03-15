package gopush_test

import (
	"testing"

	"github.com/DataWraith/gopush"
)

func TestPushingLiterals(t *testing.T) {
	interpreter := gopush.NewInterpreter(gopush.DefaultOptions)
	interpreter.Run("3 3.1415926535")

	if interpreter.Stacks["integer"].Peek().(int64) != 3 {
		t.Error("expected integer stack to contain 3")
	}

	if interpreter.Stacks["float"].Peek().(float64) != 3.1415926535 {
		t.Error("expected float stack to contain 3.1415926535")
	}
}
