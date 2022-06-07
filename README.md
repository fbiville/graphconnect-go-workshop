# Mastering Neo4j with Go and GoGM

This repository defines a set of exercises covering:

1. the basics of [Go](https://go.dev/)
2. the [Go driver](https://github.com/neo4j/neo4j-go-driver) for [Neo4j](https://neo4j.com/)
3. the [GoGM](https://github.com/mindstand/gogm) mapping framework

## Prerequisites

You need [Go 1.18+](https://go.dev/dl/) installed on your machine along with [Docker](https://docs.docker.com/get-docker/) to run the Neo4j test containers.

## Structure

Every exercise is defined as a standard Go test.
You simply have to follow the order and fill the blanks (i.e. `TODO` 
comments), one by one.

If you do not want to use an IDE, you can use the command line to run a 
specific test like this:

```shell
go test -v -run TestVariablesAndBasicTypes/'try some built-in types' ./1-golang-intro/...
```

That commands run the test named `try some built-in types`, nested in 
`TestVariablesAndBasicTypes` in the `1-golang-intro` module.