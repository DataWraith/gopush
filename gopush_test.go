package gopush_test

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
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

// Helper function to find test suites
func findTestSuites(directory string, t *testing.T) (testsuites []string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		// Check if 2-program.push exists
		if _, err := os.Stat(filepath.Join(path, "2-program.push")); err == nil {
			// There should be 3-expected.push if we have 2-program.push
			if _, err := os.Stat(filepath.Join(path, "3-expected.push")); os.IsNotExist(err) {
				t.Errorf("expected file \"3-expected.push\" in directory %q", path)
			} else {
				// We have both 2-program.push and 3-expected.push, this is a test suite!
				testsuites = append(testsuites, path)
			}
		}

		return nil
	})

	return
}

// Compares two float stacks for equality within epsilon
func compareFloatStacks(s1, s2 *gopush.Stack, epsilon float64) bool {
	if len(s1.Stack) != len(s2.Stack) {
		return false
	}

	for i := 0; i < len(s1.Stack); i++ {
		if math.Abs(s1.Stack[i].(float64)-s2.Stack[i].(float64)) > epsilon {
			return false
		}
	}

	return true
}

// This goes through the test suite under tests/ and runs every single example
func TestSuite(t *testing.T) {
	var testOptions gopush.Options
	var interpreter, expInterpreter *gopush.Interpreter

	testsuites := findTestSuites("tests", t)
	for _, ts := range testsuites {

		testOptions = gopush.DefaultOptions
		testOptions.TopLevelPopCode = true
		testOptions.RandomSeed = 1138

		var setup, program, expected []byte
		var err error

		// Open 1-setup.push (if present)
		if _, err = os.Stat(filepath.Join(ts, "1-setup.push")); err == nil {
			setup, err = ioutil.ReadFile(filepath.Join(ts, "1-setup.push"))
			if err != nil {
				t.Fatalf("error while reading %q", filepath.Join(ts, "1-setup.push"))
			}
		}

		// Open 2-program.push
		program, err = ioutil.ReadFile(filepath.Join(ts, "2-program.push"))
		if err != nil {
			t.Fatalf("error while reading %q", filepath.Join(ts, "2-program.push"))
		}

		// Open 3-expected.push
		expected, err = ioutil.ReadFile(filepath.Join(ts, "3-expected.push"))
		if err != nil {
			t.Fatalf("error while reading %q", filepath.Join(ts, "3-expected.push"))
		}

		interpreter = gopush.NewInterpreter(testOptions)

		// Run the setup program
		err = interpreter.Run(string(setup))
		if err != nil {
			t.Fatalf("error while setting up test suite %q: %v", ts, err)
		}

		// Run the test program
		err = interpreter.Run(string(program))
		if err != nil {
			t.Fatalf("error while running test suite %q: %v", ts, err)
		}

		expInterpreter = gopush.NewInterpreter(testOptions)

		// Run the expected program
		err = expInterpreter.Run(string(expected))
		if err != nil {
			t.Fatalf("error while running expected program of test suite %q: %v", ts, err)
		}

		for name, stack := range interpreter.Stacks {
			// Missing and empty stacks are equivalent
			if len(stack.Stack) == 0 && len(expInterpreter.Stacks[name].Stack) == 0 {
				continue
			}

			// We need to separately handle the FLOAT stack since
			// the calculations using floating point are giving
			// slightly different values on Drone.io
			if name == "float" {
				if !compareFloatStacks(stack, expInterpreter.Stacks[name], 1.0/1000000) {
					t.Errorf("testsuite %q: stack float does not equal expected. Expected: \n%v\n, got: \n%v\n", ts, expInterpreter.Stacks[name].Stack, stack.Stack)
				}
				continue
			}

			if !reflect.DeepEqual(stack.Stack, expInterpreter.Stacks[name].Stack) {
				t.Errorf("testsuite %q: stack %s does not equal expected. Expected: \n%v\n, got: \n%v\n", ts, name, expInterpreter.Stacks[name].Stack, stack.Stack)
			}
		}
	}
}
