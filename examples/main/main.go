package main

import (
	graphql "github.com/ichaly/go-gql"
	"log"
	"net/http"
)

func main() {
	http.Handle("/graphql", graphql.NewServer())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
