package workshop_test

import (
	"context"
	workshop "graphconnect/go-driver/pkg"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestNeo4jDriverQueryExecution(outer *testing.T) {

	ctx := context.Background()

	neo4jContainer, err := workshop.StartNeo4jContainer(ctx, workshop.ContainerConfiguration{
		Neo4jVersion: "4.4",
		Username:     username,
		Password:     password,
	})
	if err != nil {
		outer.Fatalf("Could not start container: %v", err)
	}
	defer func() {
		if err := neo4jContainer.Terminate(ctx); err != nil {
			outer.Fatalf("Could not stop container: %v", err)
		}
	}()

	driver := createDriver(outer, ctx, neo4jContainer)
	defer func() {
		if err := driver.Close(); err != nil {
			outer.Fatalf("Could not close driver: %v", err)
		}
	}()

	// Run `go test -v -run TestNeo4jDriverQueryExecution/'runs an auto-commit query' ./2-neo4j-go-driver/...`
	outer.Run("runs an auto-commit query", func(t *testing.T) {
		// an auto-commit query automatically starts a transaction on the server side
		// starting a client-side transaction for autocommit queries is forbidden and fails
		// the only auto-commit queries today are CALL {} IN TRANSACTIONS (introduced in 4.4) and USING PERIODIC COMMIT (gone in 5.0)
		// these need to be executed with Session#Run
		session := driver.NewSession(neo4j.SessionConfig{})
		defer func() {
			if err := session.Close(); err != nil {
				t.Errorf("Could not close session: %v", err)
			}
		}()
		// this is for illustration purposes only ;)
		query := "CALL { RETURN 42 AS answer } IN TRANSACTIONS RETURN answer"

		var err error
		var result neo4j.Result
		// TODO: use the correct Session method to run this autocommit query

		if err != nil {
			t.Errorf("Expected query to successfully execute but did not: %v", err)
		}
		if answer := extractAnswer(t, result); answer != 42 {
			t.Errorf("Expected 42 from %s but got: %v", query, answer)
		}
	})
	// Run `go test -v -run TestNeo4jDriverQueryExecution/'runs a read transaction, with parameters' ./2-neo4j-go-driver/...`
	outer.Run("runs a read transaction, with parameters", func(t *testing.T) {
		session := driver.NewSession(neo4j.SessionConfig{})
		defer func() {
			if err := session.Close(); err != nil {
				t.Errorf("Could not close session: %v", err)
			}
		}()
		// this defines a transaction function
		// this function may be called several times by the driver
		// ... until transient errors stop happening or the driver exhausts all attempts
		transactionFunction := func(tx neo4j.Transaction) (any, error) {
			result, err := tx.Run(
				"RETURN REDUCE(sum=0, power IN $powersOfTwo | sum+power) AS answer",
				map[string]any{
					"powersOfTwo": []int{2, 8, 32},
				})
			if err != nil {
				return nil, err
			}
			return extractAnswer(t, result), nil
		}

		var err error
		var answer any
		// TODO: remove the next line and use the correct Session method to run this read transaction
		transactionFunction(nil)

		if err != nil {
			t.Errorf("Expected query to successfully execute but did not: %v", err)
		}
		if answer != int64(42) {
			t.Errorf("Expected 42 from read transaction but got: %v", answer)
		}
	})
}

func extractAnswer(t *testing.T, result neo4j.Result) int64 {
	record, err := result.Single()
	if err != nil { // error: 0 or more than 1 result
		t.Errorf("Expected single record, but got %v", err)
	}
	rawAnswer, found := record.Get("answer")
	if !found { // error: no "answer" column
		t.Errorf(`Expected to find the "answer" column, but did not`)
	}
	// remember every integer in Cypher is mapped to an int64 in Go!
	answer, ok := rawAnswer.(int64)
	if !ok { // error: the "answer" single value is not an int64
		t.Errorf("Expected a 64-int integer, but was not: %v", answer)
	}
	return answer
}
