package main

import (
	"context"
	"github.com/mindstand/gogm/v2"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io"
	"os"
)

func main() {
	config := gogm.Config{
		Host:          os.Args[0],
		Port:          7687,
		Protocol:      "bolt",
		PoolSize:      10,
		Username:      os.Args[1],
		Password:      os.Args[2],
		IndexStrategy: gogm.IGNORE_INDEX,
		LoadStrategy:  gogm.PATH_LOAD_STRATEGY,
	}
	_gogm, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &Hello{})
	if err != nil {
		panic(err)
	}

	defer handleClose(_gogm)

	sess, err := _gogm.NewSessionV2(gogm.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "neo4j"})
	if err != nil {
		panic(err)
	}

	defer handleClose(sess)

	hello := NewHello("world")
	err = sess.SaveDepth(context.Background(), hello, 0)
	if err != nil {
		panic(err)
	}
}

func handleClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}

type Hello struct {
	gogm.BaseNode
	Quote string `gogm:"name=quote"`
}

func NewHello(quote string) *Hello {
	return &Hello{Quote: quote}
}
