// every file must declare its package
// a folder cannot contain more than 1 package (and its _test variant like
// workshop and workshop_test here).
// Package name ending with _test are only allowed for Go files ending with _test
package workshop_test

import "testing"

// this is a standard Go test
// it must be a function whose name starts with Test, that accepts a parameter of type *testing.T
func TestVariablesAndBasicTypes(outer *testing.T) {

	// tests can be nested too!
	// nested tests:
	// - are declared with Run
	// - accept a name string
	// - and a callback parameter
	// variables from the outer scope can be shared between nested tests
	outer.Run("this is an already solved exercise", func(t *testing.T) {
		// this is a long declaration form
		var aVariable int = 42
		// this is a long declaration form, with delayed initialization
		var anotherVariable int
		anotherVariable = 0x2A // 42 in hexadecimal
		// this is a long declaration form, with type inference
		var yetAnotherVariable = 0b101010 // 42 in binary
		// this is a shorthand declaration (most common)
		lastVariableWePromise := 42

		if aVariable != anotherVariable ||
			aVariable != yetAnotherVariable ||
			aVariable != lastVariableWePromise {
			t.Error("Someone has broken math")
		}
	})

	outer.Run("try some built-in types", func(t *testing.T) {
		// note for later: the Cypher type system only knows 64-bit integers and floating numbers
		// the corresponding Go types exposed by the Neo4j driver are int64 and float64
		// these are NOT equivalent to int and float (system-dependent size)

		var franceIsBacon bool
		var count int64
		var average float64
		var actorName string
		// any (Go 1.18+) is also known as interface{}
		// any is equivalent to Java's Object type
		var anythingGoes any = true

		// TODO: initialize these variables with the expected values
		// do NOT change the assertions below

		if !franceIsBacon {
			// https://knowyourmeme.com/memes/france-is-bacon
			t.Errorf("Expected France to be bacon, but was: %t", franceIsBacon)
		}
		if count != 42 {
			t.Errorf("Expected count to be 42, but was: %d", count)
		}
		if average != 0.5 {
			t.Errorf("Expected average to be 0.5 but was %f", average)
		}
		// strings are compared with == and != in Go
		if actorName != "Jane Doe" {
			// use ` to avoid having to mess with double quotes
			t.Errorf(`Expected actor to have name "Jane Doe" but was %q`, actorName)
		}
		if anythingGoes != "Long live GraphConnect!" {
			// %v is a placeholder that basically works with every value
			t.Errorf(`Expected anythingGoes to equal "Long live GraphConnect!" but was %v`, anythingGoes)
		}
	})

	outer.Run("let's play with pointers", func(t *testing.T) {
		answer := 24
		// below is a pointer to an int, a.k.a. *int, which defaults to a nil address
		var pointingToTheAnswer *int
		// the variable now stores the address of the answer variable, via the & operator
		pointingToTheAnswer = &answer

		// TODO: change the answer variable to the expected value

		// * dereferences the pointer, i.e. resolves the value at the address of the pointer
		if *pointingToTheAnswer != 42 {
			t.Errorf("Expected *pointingToTheAnswer to equal 42 but was %d", *pointingToTheAnswer)
		}
	})
}
