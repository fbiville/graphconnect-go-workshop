package workshop_test

import (
	"context"
	"fmt"
	workshop "graphconnect/go-driver/pkg"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go"
)

const username = "neo4j"
const password = "s3cr3t"

func TestNeo4jDriverConnectivity(outer *testing.T) {

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
		err := neo4jContainer.Terminate(ctx)
		if err != nil {
			outer.Fatalf("Could not stop container: %v", err)
		}
	}()
	// Run `go test -v -run TestNeo4jDriverConnectivity/'creates a Neo4j driver and verify connectivity' ./2-neo4j-go-driver/...`
	outer.Run("creates a Neo4j driver and verify connectivity", func(t *testing.T) {
		// TODO: fix the createDriver function below
		driver := createDriver(t, ctx, neo4jContainer)
		defer func() {
			if err := driver.Close(); err != nil {
				t.Fatalf("Could not close driver: %v", err)
			}
		}()

		err := driver.VerifyConnectivity()

		if err != nil {
			t.Fatalf("Expected driver to connect to the container but did not: %v", err)
		}
	})
}

func createDriver(t *testing.T, ctx context.Context, container testcontainers.Container) neo4j.Driver {
	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		t.Fatalf("Could not get mapped Bolt port: %v", err)
	}
	uri := fmt.Sprintf("neo4j://localhost:%d", port.Int())
	auth := neo4j.BasicAuth(username, password, "")
	// TODO: create the driver to connect to the running container
	panic(fmt.Errorf("connect driver to %s with %v", uri, auth))
}
