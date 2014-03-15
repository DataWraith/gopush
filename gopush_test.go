package gopush_test

import (
	"testing"

	"github.com/DataWraith/gopush"
)

// Test that literals are correctly pushed onto their respective stacks
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

// Test that the examples from the section "Simple Examples" of the Push 3.0
// specification at http://faculty.hampshire.edu/lspector/push3-description.html
// work properly
func TestSimpleExample1(t *testing.T) {
	interpreter := gopush.NewInterpreter(gopush.DefaultOptions)
	err := interpreter.Run("( 2 3 INTEGER.* 4.1 5.2 FLOAT.+ TRUE FALSE BOOLEAN.OR )")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if interpreter.Stacks["boolean"].Pop().(bool) != true {
		t.Error("expected boolean stack to contain TRUE")
	}

	if interpreter.Stacks["float"].Pop().(float64) != 9.3 {
		t.Error("expected float stack to contain 9.3")
	}

	if interpreter.Stacks["integer"].Pop().(int64) != 6 {
		t.Error("expected integer stack to contain 6")
	}
}

func TestSimpleExample2(t *testing.T) {
	interpreter := gopush.NewInterpreter(gopush.DefaultOptions)
	err := interpreter.Run("( 5 1.23 INTEGER.+ ( 4 ) INTEGER.- 5.67 FLOAT.* )")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if interpreter.Stacks["float"].Pop().(float64) != 6.9741 {
		t.Error("expected float stack to contain 6.9741")
	}

	if interpreter.Stacks["integer"].Pop().(int64) != 1 {
		t.Error("expected integer stack to contain 1")
	}
}
