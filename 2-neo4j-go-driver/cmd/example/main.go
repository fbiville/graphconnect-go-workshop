package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io"
	"os"
	"strings"
)

func main() {
	driver, err := neo4j.NewDriver(os.Args[1], neo4j.BasicAuth(os.Args[2], os.Args[3], ""))
	if err != nil {
		panic(err)
	}
	defer handleClose(driver)
	session := driver.NewSession(neo4j.SessionConfig{})
	defer handleClose(session)
	result, err := session.ReadTransaction(sayHello)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Program says: %q", result)
}

func sayHello(tx neo4j.Transaction) (any, error) {
	result, err := tx.Run(
		`RETURN reduce(acc = "", letter IN $letters | acc + letter) AS hello`,
		map[string]interface{}{
			"letters": strings.Split("Hello, GraphConnect!", ""),
		})
	if err != nil {
		return nil, err
	}
	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	rawHello, found := record.Get("hello")
	if !found {
		return nil, fmt.Errorf("expected 'hello' column in result, none found")
	}
	hello, ok := rawHello.(string)
	if !ok {
		return nil, fmt.Errorf("expected 'hello' column to be a string, but was not")
	}
	return hello, nil
}

func handleClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}
