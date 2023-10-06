package testing

import (
	"github.com/andyzhou/grit"
	"github.com/andyzhou/grit/face"
)

const (
	DBTag = "test"
)

func InitDB() (*face.DB, error) {
	g := grit.GetGrit()
	db, err := g.OpenDB(DBTag)
	return db, err
}
