package main

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	fmt.Println("test")

	repo, err := git.PlainOpen("/home/todd/Projects/ham-go")
	if err != nil {
		log.Fatalf("error opening repository: %s", err)
	}
	lo := &git.LogOptions{
		Order: git.LogOrderCommitterTime,
	}
	logs, err := repo.Log(lo)
	if err != nil {
		log.Fatalf("error retrieving logs: %s", err)
	}
	err = logs.ForEach(func(c *object.Commit) error {
		fmt.Printf("COMMIT %v\n", c)

		return nil
	})
	if err != nil {
		log.Fatalf("error enumerating logs: %s", err)
	}
}
