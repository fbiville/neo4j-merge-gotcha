package main

import (
	"context"
	"fmt"
	"github.com/fbiville/neo4j-merge-gotcha/pkg/container"
	"github.com/fbiville/neo4j-merge-gotcha/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"sync"
)

func main() {

	ctx := context.Background()
	instance, driver, err := container.StartSingleInstance(ctx, container.ContainerConfiguration{
		Neo4jVersion: "4.4",
		Username:     "neo4j",
		Password:     "s3cr3t",
		Scheme:       "neo4j",
	})
	errors.PanicOnErr(err)
	defer func() {
		errors.PanicOnErr(instance.Terminate(ctx))
	}()
	defer func() {
		errors.PanicOnErr(driver.Close())
	}()

	createUniqueConstraint(driver)

	params := map[string]interface{}{
		"name": "Jane Doe",
	}
	goRoutines := 10_000
	group := sync.WaitGroup{}
	group.Add(goRoutines)
	for i := 0; i < goRoutines; i++ {
		go func() {
			session := driver.NewSession(neo4j.SessionConfig{})
			result, err := session.Run("MERGE (:Person {name: $name})", params)
			errors.PanicOnErr(err)
			_, err = result.Consume()
			errors.PanicOnErr(err)
			errors.PanicOnErr(session.Close())
			group.Done()
		}()
	}
	group.Wait()
	session := driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		errors.PanicOnErr(session.Close())
	}()
	result, err := session.Run("MATCH (:Person {name: $name}) RETURN COUNT(*) AS count", params)
	errors.PanicOnErr(err)
	record, err := result.Single()
	errors.PanicOnErr(err)
	count, _ := record.Get("count")
	fmt.Printf("Got %d person node(s)", count)
}

func createUniqueConstraint(driver neo4j.Driver) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		errors.PanicOnErr(session.Close())
	}()
	result, err := session.Run("CREATE CONSTRAINT unique_person_name FOR (p:Person) REQUIRE p.name IS UNIQUE", nil)
	errors.PanicOnErr(err)
	_, err = result.Consume()
	errors.PanicOnErr(err)
}
