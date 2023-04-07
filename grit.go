package grit

import (
	"github.com/andyzhou/grit/face"
	"sync"
)

/*
 * Copyright (c) 2023, AndyChow <diudiu8848@gmail.com>
 * All rights reserved.
 * api face
 */

//global variable
var (
	_grit *Grit
	_gritOnce sync.Once
)

//face info
type Grit struct {
	storage *face.Storage
}

//get single instance
func GetGrit() *Grit {
	_gritOnce.Do(func() {
		_grit = NewGrit()
	})
	return _grit
}

//construct
func NewGrit() *Grit {
	this := &Grit{
		storage: face.NewStorage(),
	}
	return this
}

//get dbs
func (f *Grit) GetDBs() []string {
	return f.storage.GetDBs()
}

//close db
func (f *Grit) CloseDB(tag string) error {
	return f.storage.CloseDB(tag)
}

//open db
func (f *Grit) OpenDB(tag string) (*face.DB, error) {
	return f.storage.OpenDB(tag)
}

//set db path
func (f *Grit) SetDBPath(path string) error {
	return f.storage.SetDBPath(path)
}