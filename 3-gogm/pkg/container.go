package workshop_gogm

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartNeo4jContainer(ctx context.Context, config ContainerConfiguration) (testcontainers.Container, error) {
	version := config.Neo4jVersion
	request := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("neo4j:%s", version),
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{
			"NEO4J_AUTH":                     config.neo4jAuthEnvVar(),
			"NEO4J_ACCEPT_LICENSE_AGREEMENT": "yes",
		},
		WaitingFor: wait.ForLog("Bolt enabled"),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}
	return container, err
}

type ContainerConfiguration struct {
	Neo4jVersion string
	Username     string
	Password     string
}

func (config ContainerConfiguration) neo4jAuthEnvVar() string {
	return fmt.Sprintf("%s/%s", config.Username, config.Password)
}

func (config ContainerConfiguration) neo4jAuthToken() neo4j.AuthToken {
	return neo4j.BasicAuth(config.Username, config.Password, "")
}
