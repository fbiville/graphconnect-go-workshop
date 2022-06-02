package workshop_gogm

import (
	"context"
	"github.com/mindstand/gogm/v2"
	"testing"
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
			outer.Errorf("Could not start container: %v", err)
		}
	}()
	outer.Run("Create a GoGM instance", func(t *testing.T) {
		containerIP, err := neo4jContainer.ContainerIP(ctx)
		if err != nil {
			t.Fatalf(err.Error())
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

		// gogm.New() takes in the conifg object, the primary key strategy and a list of nodes it will use
		// anyone can create a custom primary key strategy. For this example we will use UUID
		// the default option just uses neo4j's built in int64 graph id
		// all nodes and Edge interface implementations must be provided here
		// if gogm recieves a node it doesn't recognize it will error out
		// this is intended to reduce the amount of reflect calls at runtime
		_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &Project{}, &Person{}, &Topic{}, &WorksOnEdge{})
		if err != nil {
			panic(err)
		}
		defer _gogm.Close()

		// you can also set a global gogm to be accessed via gogm.G()
		gogm.SetGlobalGogm(_gogm)
	})
}
