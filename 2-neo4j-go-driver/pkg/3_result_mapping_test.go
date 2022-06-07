package workshop_test

import (
	"context"
	"fmt"
	workshop "graphconnect/go-driver/pkg"
	"reflect"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestNeo4jDriverResultMapping(outer *testing.T) {

	ctx := context.Background()

	neo4jContainer, err := workshop.StartNeo4jContainer(ctx, workshop.ContainerConfiguration{
		Neo4jVersion: "4.4",
		Username:     username,
		Password:     password,
	})
	if err != nil {
		outer.Errorf("Could not start container: %v", err)
	}
	defer func() {
		if err := neo4jContainer.Terminate(ctx); err != nil {
			outer.Errorf("Could not stop container: %v", err)
		}
	}()

	driver := createDriver(outer, ctx, neo4jContainer)
	defer func() {
		if err := driver.Close(); err != nil {
			outer.Errorf("Could not close driver: %v", err)
		}
	}()
	insertSmallGraph(outer, driver)
	// Run `go test -v -run TestNeo4jDriverResultMapping/'extracts persons working on projects' ./2-neo4j-go-driver/...`
	outer.Run("extracts persons working on projects", func(t *testing.T) {
		session := driver.NewSession(neo4j.SessionConfig{})
		defer func() {
			if err := session.Close(); err != nil {
				t.Errorf("Session could not close: %v", err)
			}
		}()

		names, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			var records []*neo4j.Record
			query := "MATCH (p:Person)-[:WORKS_ON]->(:Project) RETURN p ORDER BY p.name ASC"
			// TODO: remove next line and run query + collect records
			fmt.Println(query)

			names := make([]string, len(records))
			for i, record := range records {
				var name string
				// TODO: remove next line and extract the name property from the returned nodes
				fmt.Println(record)
				names[i] = name
			}
			return names, nil
		})

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		expected := []string{"Eric", "Florent", "Nikita"}
		if !reflect.DeepEqual(names, expected) {
			t.Errorf("Expected %v, got: %v", expected, names)
		}
	})
	// Run `go test -v -run TestNeo4jDriverResultMapping/'extracts projects, sorted by maintainer count' ./2-neo4j-go-driver/...`
	outer.Run("extracts projects, sorted by maintainer count", func(t *testing.T) {
		session := driver.NewSession(neo4j.SessionConfig{})
		defer func() {
			if err := session.Close(); err != nil {
				t.Errorf("Session could not close: %v", err)
			}
		}()

		names, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			var result neo4j.Result
			query := "MATCH (p:Project)<-[:WORKS_ON]-(pe:Person) WITH p, size(collect(pe)) AS count RETURN p ORDER BY count DESC"
			// TODO: remove next line and run query + iterate over results
			fmt.Println(query)

			var names []string
			for result.Next() {
				record := result.Record()
				// TODO: remove next line and extract project name + append it to names slice
				fmt.Println(record)
			}
			err := result.Err()
			if err != nil {
				return nil, err
			}
			return names, nil
		})

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		expected := []string{"GoGM", "Go Driver"}
		if !reflect.DeepEqual(names, expected) {
			t.Errorf("Expected %v, got: %v", expected, names)
		}
	})
}

func insertSmallGraph(t *testing.T, driver neo4j.Driver) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(); err != nil {
			t.Fatalf("Session could not close: %v", err)
		}
	}()

	// note: with Neo4j, you cannot mix schema and data operations
	// ... hence the two separate transactions below
	if _, err := session.WriteTransaction(createIndices); err != nil {
		t.Fatalf("Could not create indices: %v", err)
	}
	if _, err := session.WriteTransaction(insertData); err != nil {
		t.Fatalf("Could not create data: %v", err)
	}
}

func createIndices(tx neo4j.Transaction) (any, error) {
	queries := []string{
		"CREATE INDEX FOR (t:Topic) ON (t.name)",
		"CREATE INDEX FOR (pe:Person) ON (pe.name)",
		"CREATE INDEX FOR (p:Project) ON (p.name)",
	}
	for _, query := range queries {
		result, err := tx.Run(query, nil)
		if err != nil {
			return nil, err
		}
		if _, err := result.Consume(); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func insertData(tx neo4j.Transaction) (any, error) {
	result, err := tx.Run(`
MERGE (neo4j:Topic {name: "Neo4j"})
MERGE (goDriver:Project {name: "Go Driver"})
MERGE (gogm:Project {name: "GoGM"})
MERGE (album:MusicProject {name: "TBD"})
MERGE (eric:Person {name: "Eric"})
MERGE (nikita:Person {name: "Nikita"})
MERGE (florent:Person {name: "Florent"})
MERGE (john:Person {name: "John"})
MERGE (gogm)-[:RELATES_TO]->(neo4j)
MERGE (eric)-[:WORKS_ON]->(gogm)
MERGE (nikita)-[:WORKS_ON]->(gogm)
MERGE (goDriver)-[:RELATES_TO]->(neo4j)
MERGE (florent)-[:WORKS_ON]->(goDriver)
MERGE (john)-[:WORKS_ON]->(album)
`, nil)
	if err != nil {
		return nil, err
	}
	summary, err := result.Consume()
	if err != nil {
		return nil, err
	}
	return summary, nil
}
