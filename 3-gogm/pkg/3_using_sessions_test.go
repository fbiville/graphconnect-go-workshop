package workshop_gogm

import (
	"context"
	"fmt"
	"github.com/mindstand/gogm/v2"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

func TestUseSessions(outer *testing.T) {
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
		err = initGogm(ctx, neo4jContainer)
		if err != nil {
			t.Fatal(err.Error())
		}

		// assemble the graph in memory
		neo4jTopic := NewTopic("neo4j")

		eric := NewPerson("Eric")
		nikita := NewPerson("Nikita")
		florent := NewPerson("Florent")

		gogmProject := NewProject("gogm", "software")
		goDriverProject := NewProject("neo4j-go-driver", "software")

		// first lets relate the topics. There is a cli to generate linking functions, but for this example we'll do them manually
		// the important thing to note is that they must be related on both sides
		neo4jTopic.Projects = []*Project{gogmProject, goDriverProject}
		gogmProject.Topics = []*Topic{neo4jTopic}
		goDriverProject.Topics = []*Topic{neo4jTopic}

		// now let's create the special edges that store the role each person has to their project
		ericToGogm := NewWorksOnEdge(eric, gogmProject, "Lead")
		nikitaToGogm := NewWorksOnEdge(nikita, gogmProject, "Lead")
		florentToDriver := NewWorksOnEdge(florent, goDriverProject, "Lead")

		// now lets set up the relationship
		gogmProject.People = []*WorksOnEdge{ericToGogm, nikitaToGogm}
		eric.Projects = []*WorksOnEdge{ericToGogm}
		nikita.Projects = []*WorksOnEdge{nikitaToGogm}
		goDriverProject.People = []*WorksOnEdge{florentToDriver}
		florent.Projects = []*WorksOnEdge{florentToDriver}

		// now let's create a session so that we can save this
		// we're accessing the session through gogm.G() since we set in initGoGM
		// always use NewSessionV2 since new session deprecated
		// session config is an alias of the driver session config
		// we'll point at the default database (neo4j) and use read write
		sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{
			AccessMode:   neo4j.AccessModeWrite,
			DatabaseName: "neo4j",
		})
		if err != nil {
			t.Fatal(err.Error())
		}

		// always make sure to defer closing the session
		defer func() {
			err = sess.Close()
			if err != nil {
				t.Fatal(err.Error())
			}
		}()

		// now we can save
		// there are a few options
		// we can save the topic at a depth of 2, or we can save each project at a depth of 1
		// for simplicity we'll do topic at a depth of 2
		// Save always expects a pointer and will return an error if not given one
		err = sess.SaveDepth(ctx, neo4jTopic, 2)
		if err != nil {
			t.Fatal(err.Error())
		}

		// under the hood save depth runs in a transaction, you can also explicitly run it in a transaction with the following code
		//err = sess.ManagedTransaction(ctx, func(tx gogm.TransactionV2) error {
		//	return tx.SaveDepth(ctx, neo4jTopic, 2)
		//})
		//if err != nil {
		//	t.Fatal(err)
		//}
		// these 2 approaches essentially do the same thing

		// when gogm saves it is updating primary keys if the node is new and is detecting if relationships were added/deleted
		// you will notice LoadMap as a member of all nodes. Gogm uses this to track the state of the relationships between saves and loads

		// the following is an example illustrating this
		var loadedNeo4jTopic Topic
		err = sess.LoadDepth(ctx, &loadedNeo4jTopic, neo4jTopic.UUID, 1)
		if err != nil {
			t.Fatal(err)
		}

		// now we can check if the relationships were loaded correctly
		var loadedGogm *Project
		var gogmIndex int
		for i, project := range loadedNeo4jTopic.Projects {
			if project.Name == gogmProject.Name && project.UUID == gogmProject.UUID {
				loadedGogm = project
				gogmIndex = i
				continue
			} else if project.Name == goDriverProject.Name && project.UUID == goDriverProject.UUID {
				continue
			} else {
				t.Fatal("did not load edges correctly")
			}
		}

		// now we can illustrate removing a relationship
		// wipe the gogm topics map
		loadedGogm.Topics = []*Topic{}
		// recreated the projects slice with just the driver
		loadedNeo4jTopic.Projects = remove(loadedNeo4jTopic.Projects, gogmIndex)

		// now we can save the topic at a depth of 1
		err = sess.SaveDepth(ctx, &loadedNeo4jTopic, 1)
		if err != nil {
			t.Fatal(err.Error())
		}

		// now we can load the loaded struct back in 1 more time to see if the relationships reflect
		var loadedNeo4jTopic2 Topic
		err = sess.LoadDepth(ctx, &loadedNeo4jTopic2, neo4jTopic.UUID, 1)
		if err != nil {
			t.Fatal(err)
		}

		// the length of projects should only be one now
		if len(loadedNeo4jTopic2.Projects) != 1 {
			t.Fatal("length of projects was not equal to 1")
		}

		// and we can verify its the driver
		if loadedNeo4jTopic2.Projects[0].UUID != goDriverProject.UUID {
			t.Fatal("loaded project was not the go driver project")
		}

		// that's the basic usage of GoGM. Since we deferred closing the session we don't have to do anything else here to close it
	})
}

func initGogm(ctx context.Context, neo4jContainer testcontainers.Container) error {
	containerIP, err := neo4jContainer.ContainerIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container IP: %w", err)
	}
	config := gogm.Config{
		Host:          containerIP,
		Port:          7687,
		Protocol:      "bolt",
		PoolSize:      10,
		Username:      username,
		Password:      password,
		IndexStrategy: gogm.ASSERT_INDEX,
		LoadStrategy:  gogm.PATH_LOAD_STRATEGY,
	}

	_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &Project{}, &Person{}, &Topic{}, &WorksOnEdge{})
	if err != nil {
		return fmt.Errorf("failed to init gogm: %w", err)
	}

	gogm.SetGlobalGogm(_gogm)
	return nil
}

func remove[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
