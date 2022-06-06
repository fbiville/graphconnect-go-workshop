package workshop_gogm

import (
	"context"
	"testing"

	"github.com/mindstand/gogm/v2"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
			outer.Errorf("Could not stop container: %v", err)
		}
	}()

	if err = initGogm(ctx, neo4jContainer); err != nil {
		outer.Fatal("error initializing gogm, did you finish 1_defining_a_schema?", err.Error())
	}

	outer.Run("save data using gogm", func(t *testing.T) {
		// assemble the graph in memory
		neo4jTopic := &Topic{Name: "neo4j"}

		eric := &Person{Name: "Eric"}
		nikita := &Person{Name: "Nikita"}
		florent := &Person{Name: "Florent"}

		gogmProject := &Project{Name: "gogm", Type: "software"}
		goDriverProject := &Project{Name: "neo4j-go-driver", Type: "software"}

		// first lets relate the topics. There is a cli to generate linking functions, but for this example we'll do them manually
		// the important thing to note is that they must be related on both sides
		neo4jTopic.Projects = []*Project{gogmProject, goDriverProject}
		gogmProject.Topics = []*Topic{neo4jTopic}
		goDriverProject.Topics = []*Topic{neo4jTopic}

		// now let's create the special edges that store the role each person has to their project
		ericToGogm := &WorksOnEdge{Start: eric, End: gogmProject, Role: "Lead"}
		nikitaToGogm := &WorksOnEdge{Start: nikita, End: gogmProject, Role: "Lead"}
		florentToDriver := &WorksOnEdge{Start: florent, End: goDriverProject, Role: "Lead"}

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
			if err = sess.Close(); err != nil {
				t.Fatal(err.Error())
			}
		}()

		// now we can save
		// there are a few options
		// we can save the topic at a depth of 2, or we can save each project at a depth of 1
		// for simplicity we'll do topic at a depth of 2
		// Save always expects a pointer and will return an error if not given one
		if err = sess.SaveDepth(ctx, neo4jTopic, 2); err != nil {
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

		// Since we deferred closing the session we don't have to do anything else here to close it
	})

	outer.Run("loading and changing data using gogm", func(t *testing.T) {
		// create the session
		sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{
			AccessMode:   neo4j.AccessModeWrite,
			DatabaseName: "neo4j",
		})
		if err != nil {
			t.Fatal(err.Error())
		}

		// always make sure to defer closing the session
		defer func() {
			if err = sess.Close(); err != nil {
				t.Fatal(err.Error())
			}
		}()

		// make sure nothing's left from the previous test
		if err = clearDb(ctx, sess); err != nil {
			t.Fatal("Failed to clear database, error:", err)
		}
		// insert our test data
		if err = insertData(ctx, sess); err != nil {
			t.Fatal("Failed to insert data, error:", err)
		}

		// when gogm saves it is updating primary keys if the node is new and is detecting if relationships were added/deleted
		// you will notice LoadMap as a member of all nodes. Gogm uses this to track the state of the relationships between saves and loads

		// the following is an example illustrating this
		var loadedProject Project
		err = sess.Query(ctx, `match p=(n:Project{name:$name})-[]-(:Person) return p`, map[string]interface{}{
			"name": "GoGM",
		}, &loadedProject)
		if err != nil {
			t.Fatal("Failed to query the project, error:", err)
		}

		if len(loadedProject.People) != 2 {
			t.Fatal("No people were loaded in query")
		}

		// now we can illustrate removing a relationship
		// TODO: wipe the people list in loadedProject

		// now we can save the project at a depth of 1
		if err = sess.SaveDepth(ctx, &loadedProject, 2); err != nil {
			t.Fatal(err.Error())
		}

		// now we can load the loaded struct back in 1 more time to see if the relationships reflect
		var loadedProject2 Topic

		if err = sess.LoadDepth(ctx, &loadedProject2, loadedProject.UUID, 1); err != nil {
			t.Fatal(err)
		}

		// the length of projects should only be one now
		if len(loadedProject2.Projects) == 0 {
			t.Error("length of projects was not equal to 0")
		}
	})

	outer.Run("closing remarks", func(t *testing.T) {
		// create the session
		sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{
			AccessMode:   neo4j.AccessModeWrite,
			DatabaseName: "neo4j",
		})
		if err != nil {
			t.Fatal(err.Error())
		}

		// always make sure to defer closing the session
		defer func() {
			if err = sess.Close(); err != nil {
				t.Fatal(err.Error())
			}
		}()

		// make sure nothing's left from the previous test
		if err = clearDb(ctx, sess); err != nil {
			t.Fatal("Failed to clear database, error:", err)
		}
		// insert our test data
		if err = insertData(ctx, sess); err != nil {
			t.Fatal("Failed to insert data, error:", err)
		}

		// a very hacky query that avoids using APOC ;)
		query := `MATCH (p:Person)
		WITH COLLECT(p.name) as names
		WITH REDUCE(merged= "" , name IN names | merged + name + ' and ') as joined
		RETURN LEFT(joined, SIZE(joined) - 5) + " thank you for joining the " + $conference + " Go Workshop!"`

		res, _, err := sess.QueryRaw(ctx, query, map[string]interface{}{
			// TODO: Oops, I forgot to se the conference name :)
			"conference": "GraphConnect",
		})
		if err != nil {
			t.Errorf("Did you forget to set $conference?")
		}

		resStr, ok := res[0][0].(string)
		if !ok {
			t.Errorf("Couldn't decode resStr")
		}

		t.Logf("Query response: %s", resStr)
	})
}

// Helper that delete everything in the database
func clearDb(ctx context.Context, sess gogm.SessionV2) error {
	_, _, err := sess.QueryRaw(ctx, "match(n) detach delete n", map[string]interface{}{})
	return err
}

// Helper that inserts a nice bit of test data
// Note that in a real scenario GoGM will generate link functions and remove most of this code
func insertData(ctx context.Context, sess gogm.SessionV2) error {
	// Florent subgraph
	florent := &Person{Name: "Florent"}
	goDriverProject := &Project{Name: "Go Driver"}

	florent2goDriver := &WorksOnEdge{Start: florent, End: goDriverProject}
	florent.Projects = []*WorksOnEdge{florent2goDriver}
	goDriverProject.People = []*WorksOnEdge{florent2goDriver}

	// Mindstand subgraph
	nikita := &Person{Name: "Nikita"}
	eric := &Person{Name: "Eric"}
	gogmProject := &Project{Name: "GoGM"}

	nikita2gogm := &WorksOnEdge{Start: nikita, End: gogmProject, Role: "Lead"}
	nikita.Projects = append(nikita.Projects, nikita2gogm)
	gogmProject.People = append(gogmProject.People, nikita2gogm)

	eric2gogm := &WorksOnEdge{Start: eric, End: gogmProject, Role: "Lead"}
	eric.Projects = append(eric.Projects, eric2gogm)
	gogmProject.People = append(gogmProject.People, eric2gogm)

	// the common neo4jTopic
	neo4jTopic := &Topic{Name: "neo4j"}
	neo4jTopic.Projects = []*Project{gogmProject, goDriverProject}
	gogmProject.Topics = []*Topic{neo4jTopic}
	goDriverProject.Topics = []*Topic{neo4jTopic}

	// save the entire graph
	return sess.SaveDepth(ctx, neo4jTopic, 2)
}
