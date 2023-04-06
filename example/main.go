package main

import (
	"fmt"
	"github.com/andyzhou/grit"
	"os"
)

func main() {
	dbTag := "test"
	g := grit.GetGrit()
	db, err := g.OpenDB(dbTag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.CloseDB()
}
