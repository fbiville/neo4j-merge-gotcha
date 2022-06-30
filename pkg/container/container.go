package container

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerConfiguration struct {
	Neo4jVersion string
	Username     string
	Password     string
	Scheme       string
}

func (config ContainerConfiguration) neo4jAuthEnvVar() string {
	return fmt.Sprintf("%s/%s", config.Username, config.Password)
}

func (config ContainerConfiguration) neo4jAuthToken() neo4j.AuthToken {
	return neo4j.BasicAuth(config.Username, config.Password, "")
}

func StartSingleInstance(ctx context.Context, config ContainerConfiguration) (testcontainers.Container, neo4j.Driver, error) {
	version := config.Neo4jVersion
	request := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("neo4j:%s", version),
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{
			"NEO4J_AUTH":                     config.neo4jAuthEnvVar(),
			"NEO4J_ACCEPT_LICENSE_AGREEMENT": "yes",
		},
		WaitingFor: boltReadyStrategy(),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}
	driver, err := newNeo4jDriver(ctx, config.Scheme, container, config.neo4jAuthToken())
	return container, driver, err
}

func boltReadyStrategy() *wait.LogStrategy {
	return wait.ForLog("Bolt enabled")
}

func newNeo4jDriver(ctx context.Context, scheme string, container testcontainers.Container, auth neo4j.AuthToken) (neo4j.Driver, error) {
	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		return nil, err
	}
	return newDriver(scheme, port.Int(), auth)
}

func newDriver(scheme string, port int, auth neo4j.AuthToken) (neo4j.Driver, error) {
	uri := fmt.Sprintf("%s://localhost:%d", scheme, port)
	return neo4j.NewDriver(uri, auth)
}
