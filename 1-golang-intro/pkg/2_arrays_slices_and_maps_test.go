package workshop_test

import (
	"reflect"
	"strings"
	"testing"
)

func TestArraysSlicesAndMaps(outer *testing.T) {

	outer.Run("quick look at arrays", func(t *testing.T) {
		// arrays have a fixed size, which is a part of their type
		// [3]int is not the same as [5]int
		values := [3]int{1, 2, 3}

		// solution
		values[0] = 321

		if values[0] != 321 {
			t.Errorf("expected first element to equal 321 but was %d", values[0])
		}
	})

	outer.Run("who wants a slice of that?", func(t *testing.T) {
		// slices are dynamically-growing lists
		// they resemble arrays except they don't declare their size
		// []int is not the same as [3]int
		values := []int{1, 2, 3}

		// solution
		values = append(values, 42)

		// hint: since the slice sometimes need to grow, it needs to be reallocated and its address will change!
		// to cope with this problem, make sure to re-assign `values` to the result of append
		// mySlice := append(mySlice, someValue) is a very common pattern in Go
		// bonus: in this example, mySlice can be nil (as in: `var mySlice []SomeType`) and append will work properly

		if len(values) != 4 {
			t.Errorf("expected slice to have size 4 but was %d", len(values))
		}
		if values[3] != 42 {
			t.Errorf("expected fourth element to equal 42 but was %d", values[3])
		}
	})

	outer.Run("let's iterate", func(t *testing.T) {
		// let's pre-allocate a slice of 3 elements
		values := make([]int, 3)
		// any kind of loops in Go is with the keyword "for" (no while, do-while)
		for i := 0; i < 3; i++ {
			values[i] = i + 1
		}
		var doubledValues []int

		for _, curVal := range values {
			doubledValues = append(doubledValues, curVal*2)
		}
		// hint: range produces 1 to 2 values: the current index and the current value
		// 		 you can ignore the index (and in general any value/parameter/...) with _
		// 		 e.g.: for _, currentValue := range someSlice { ... }
		// 		 use append again to insert elements into `doubledValues`

		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(doubledValues, expected) {
			t.Errorf("expected doubled values to be %v, but was: %v", expected, doubledValues)
		}
	})

	outer.Run("let's iterate over maps!", func(t *testing.T) {
		// maps (or dictionaries) can be declared as such
		digits := map[string]int{
			"One":   1,
			"Two":   2,
			"Three": 3,
			"Four":  4,
		}
		var song strings.Builder

		// TODO: understand and fix the issue (it may take more than 1 test run to realize the failure)
		// hint: there is a function below that might help!
		for _, key := range sortedKeys(digits) {
			song.WriteString(key)
			song.WriteString(" ")
		}

		song.WriteString("by Feist")
		songName := song.String()
		if songName != "One Two Three Four by Feist" {
			t.Errorf(`expected song name to be "One Two Three Four by Feist", but was: %q`, songName)
		}
	})
}

func sortedKeys(digits map[string]int) []string {
	result := make([]string, len(digits))
	for key, value := range digits {
		result[value-1] = key
	}
	return result
}
