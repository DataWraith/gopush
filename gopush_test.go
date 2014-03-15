package gopush_test

import (
	"testing"

	"github.com/DataWraith/gopush"
)

func TestPushingLiterals(t *testing.T) {
	interpreter := gopush.NewInterpreter(gopush.DefaultOptions)
	err := interpreter.Run("3 3.1415926535 FALSE TRUE")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if interpreter.Stacks["integer"].Pop().(int64) != 3 {
		t.Error("expected integer stack to contain 3")
	}

	if interpreter.Stacks["float"].Pop().(float64) != 3.1415926535 {
		t.Error("expected float stack to contain 3.1415926535")
	}

	b1 := interpreter.Stacks["boolean"].Pop().(bool)
	b2 := interpreter.Stacks["boolean"].Pop().(bool)

	if b1 != true {
		t.Error("expected top of the boolean stack to contain TRUE")
	}

	if b2 != false {
		t.Error("expected bottom of the boolean stack to contain FALSE")
	}
}
