package workshop_gogm

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/mindstand/gogm/v2"
	"github.com/testcontainers/testcontainers-go"
)

func TestInitializeGogm(outer *testing.T) {
	ctx := context.Background()
	neo4jContainer, err := StartNeo4jContainer(ctx, ContainerConfiguration{
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
	// Run `go test -v -run TestInitializeGogm/'create a gogm instance' ./3-gogm/...`
	outer.Run("create a gogm instance", func(t *testing.T) {
		// TODO: below this test, the Person, Topic, and Project have problems, fix them so that GoGM initializes

		// Checking the internal properties of WorksOnEdge
		worksOnEdgeType := reflect.TypeOf(WorksOnEdge{})
		worksOnEdge_Role, ok := worksOnEdgeType.FieldByName("Role")
		if !ok {
			t.Error("Role field does not exist on WorksOnEdge, did you rename it?")
		}
		if worksOnEdge_RoleTag := worksOnEdge_Role.Tag.Get("gogm"); worksOnEdge_RoleTag != "name=role" {
			t.Errorf("Role field's struct tag (%v) is not set correctly", worksOnEdge_RoleTag)
		}

		// Initialize GoGM, which will validate the structs defined below
		if err := initGogm(ctx, neo4jContainer); err != nil {
			t.Errorf("gogm init failed, did you fix the broken schema? Error: %v", err)
		}
	})
}

// Package schema contains the gogm schema adapted from 3_result_mapping in section 2

// Person defines a person and their relationships for the schema
type Person struct {
	// TODO: Each node must embed gogm.BaseUUIDNode (nodes can also embed gogm.BaseNode)

	// Members
	Name string `gogm:"name=name;index"`
	// Relationships
	Projects []*WorksOnEdge `gogm:"direction=outgoing;relationship=WORKS_ON"`
}

// Topic defines a topic and their relationships for the schema
type Topic struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	// Relationships
	Projects []*Project `gogm:"direction=incoming;relationship=RELATES_TO"`
}

// Project defines a project and its relationships
type Project struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	Type string `gogm:"name=project_type"`
	// Relationships
	People []*WorksOnEdge `gogm:"direction=incoming;relationship=WORKS_ON"`
	// TODO: GoGM needs a tag set to assign this struct to the other side of the `RELATES_TO` relationship
	Topics []*Topic
}

// WorksOnEdge implements gogm.Edge
// This will be phased out in a near future update for struct tags
type WorksOnEdge struct {
	gogm.BaseUUIDNode
	Start *Person
	End   *Project
	// TODO: Eric was working on this code very late at night and left an error :P
	Role string `gogm:"name=name"`
}

// These methods are needed to implement the Edge interface

func (w *WorksOnEdge) GetStartNode() interface{} {
	return w.Start
}

func (w *WorksOnEdge) GetStartNodeType() reflect.Type {
	return reflect.TypeOf(&Person{})
}

func (w *WorksOnEdge) SetStartNode(v interface{}) error {
	_start, ok := v.(*Person)
	if !ok {
		return fmt.Errorf("could not cast %T to *Person", v)
	}
	w.Start = _start
	return nil
}

func (w *WorksOnEdge) GetEndNode() interface{} {
	return w.End
}

func (w *WorksOnEdge) GetEndNodeType() reflect.Type {
	return reflect.TypeOf(&Project{})
}

func (w *WorksOnEdge) SetEndNode(v interface{}) error {
	_end, ok := v.(*Project)
	if !ok {
		return fmt.Errorf("could not cast %T to *Project", v)
	}
	w.End = _end
	return nil
}

// a little helper function that initializes gogm with the neo4j container
func initGogm(ctx context.Context, neo4jContainer testcontainers.Container) error {
	containerIP, err := neo4jContainer.ContainerIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container IP: %w", err)
	}
	config := gogm.Config{
		Host:     containerIP,
		Port:     7687,
		Protocol: "bolt",
		// use a better method to determine this
		// 10 is arbitrary
		PoolSize: 10,
		Username: username,
		Password: password,
		// INDEX strategy tells gogm whether to validate indexes exist, create them or skip verification all together
		// For this example we will use ASSERT_INDEX which will remove existing indexes for the nodes and recreate them
		IndexStrategy: gogm.ASSERT_INDEX,
		// this will generate queries based on the path, the other option is to
		// use gogm.SCHEMA_LOAD_STRATEGY which will generate queries based on the schema
		LoadStrategy: gogm.PATH_LOAD_STRATEGY,
	}

	// gogm.New() takes in the config object, the primary key strategy and a list of nodes it will use
	// anyone can create a custom primary key strategy. For this example we will use UUID
	// the default option just uses neo4j's built in int64 graph id
	// all nodes and Edge interface implementations must be provided here
	// if gogm receives a node it doesn't recognize it will error out
	// this is intended to reduce the amount of reflect calls at runtime
	_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &Project{}, &Person{}, &Topic{}, &WorksOnEdge{})
	if err != nil {
		return fmt.Errorf("failed to init gogm: %w", err)
	}

	// you can also set a global gogm to be accessed via gogm.G()
	gogm.SetGlobalGogm(_gogm)
	return nil
}
