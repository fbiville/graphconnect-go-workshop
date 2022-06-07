package workshop_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStructsAndInterfaces(outer *testing.T) {
	// Run `go test -v -run TestStructsAndInterfaces/'love-struck with structs' ./1-golang-intro/...`
	outer.Run("love-struck with structs", func(t *testing.T) {
		// TODO: first, fix the struct declaration
		type person struct {
			name    int
			address bool
		}

		// TODO: second, instantiate the struct with expected values
		eric := person{
			name:    42,
			address: true,
		}

		// reflection is used instead of == or !=
		// since int/bool and string cannot be compared directly
		if !reflect.DeepEqual(eric.name, "Eric") {
			t.Errorf(`Expected "Eric", got %v`, eric.name)
		}
		if !reflect.DeepEqual(eric.address, "Champs-Elysees, Paris FRANCE") {
			t.Errorf(`Expected "Champs-Elysees, Paris FRANCE", got %v`, eric.name)
		}
	})
	// Run `go test -v -run TestStructsAndInterfaces/'guten Tag(s)' ./1-golang-intro/...`
	outer.Run("guten Tag(s)", func(t *testing.T) {
		// TODO: add the missing tag
		type speaker struct {
			Nom                  string `json:"name"`
			Enterprise           string
			isThisGonnaBeIgnored bool
		}
		// the json package supports pointers too!
		nikita := &speaker{
			Nom:        "Nikita",
			Enterprise: "Mindstand",
		}

		result := encodeJson(t, nikita)

		if result != `{"name":"Nikita","company":"Mindstand"}`+"\n" {
			t.Errorf("Expected correct JSON to be marshaled, got: %s", result)
		}
	})
	// Run `go test -v -run TestStructsAndInterfaces/'a glimpse at interfaces' ./1-golang-intro/...`
	outer.Run("a glimpse at interfaces", func(t *testing.T) {
		// TODO: fix Neo4jChauffeur below to return the expected value
		// note: the cast here is not necessary. It would only fail at compile-time...
		// ...should the type *neo4jChauffeur not adhere to fakeDriver contract anymore
		chauffeur1 := fakeDriver(&neo4jChauffeur{defaultDb: "neo4j"})
		chauffeur2 := fakeDriver(&neo4jChauffeur{defaultDb: "custom"})

		defaultDb1 := chauffeur1.GetDefaultDatabaseName()
		if defaultDb1 != "neo4j" {
			t.Errorf(`Expected "neo4j" default DB, got %v`, defaultDb1)
		}
		defaultDb2 := chauffeur2.GetDefaultDatabaseName()
		if defaultDb2 != "custom" {
			t.Errorf(`Expected "neo4j" default DB, got %v`, defaultDb2)
		}
	})
}

type fakeDriver interface {
	GetDefaultDatabaseName() string
}

type neo4jChauffeur struct {
	defaultDb string
}

// this function has a receiver of type "pointer of Neo4jChauffeur"
// such functions are usually called methods
// this is the static equivalent of JavaScript's bind function
func (chauffeur *neo4jChauffeur) GetDefaultDatabaseName() string {
	// panic raises an error (like Java's throw)
	// unlike Java's throw, this is very rarely used!
	// indeed, errors can be returned as any other values in Go
	panic("implement me")
}

func encodeJson(t *testing.T, value any) string {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(value)
	if err != nil {
		t.Errorf("Expected no marshaling error, got %v", err)
	}
	return buffer.String()
}
