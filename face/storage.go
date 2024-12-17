package face

import (
	"errors"
	"github.com/andyzhou/grit/define"
	"sync"
)

/*
 * @author Andy Chow <diudiu8848@163.com>
 * data storage face
 * - base one leveldb
 */

//face info
type Storage struct {
	rootPath string
	dbMap    sync.Map //tag -> *DB
	Base
}

//construct
func NewStorage() *Storage {
	//self init
	this := &Storage{
		rootPath: define.DefaultDBRootPath,
		dbMap: sync.Map{},
	}
	return this
}

//quit
func (f *Storage) Quit() {
	sf := func(k, v interface{}) bool {
		db, ok := v.(*DB)
		if ok && db != nil {
			db.CloseDB()
		}
		return true
	}
	f.dbMap.Range(sf)
}

//get dbs
func (f *Storage) GetDBs() []string {
	result := make([]string, 0)
	sf := func(k, v interface{}) bool {
		if tag, ok := k.(string); ok {
			result = append(result, tag)
		}
		return true
	}
	f.dbMap.Range(sf)
	return result
}

//close db
func (f *Storage) CloseDB(tag string) error {
	//check
	if tag == "" {
		return errors.New("invalid parameter")
	}

	//hit cache first
	db := f.getDbByTag(tag)
	if db == nil {
		return errors.New("no db obj by tag")
	}

	//try close db
	err := db.CloseDB()
	if err != nil {
		return err
	}

	//remove from map
	f.dbMap.Delete(tag)
	return nil
}

//get db
func (f *Storage) GetDB(tag string) (*DB, error) {
	//check
	if tag == "" {
		return nil, errors.New("invalid parameter")
	}

	//hit cache first
	db := f.getDbByTag(tag)
	if db == nil {
		return nil, errors.New("no such db tag")
	}
	return db, nil
}

//open db
func (f *Storage) OpenDB(tag string) (*DB, error) {
	//check
	if tag == "" {
		return nil, errors.New("invalid parameter")
	}

	//hit cache first
	db := f.getDbByTag(tag)
	if db != nil {
		return db, nil
	}

	//try open db
	db = NewDB(tag)
	db.setRootPath(f.rootPath)
	err := db.OpenDB()
	if err != nil {
		return nil, err
	}

	//sync into map
	f.dbMap.Store(tag, db)
	return db, nil
}

//set db path
func (f *Storage) SetDBPath(path string) error {
	if path == "" {
		return errors.New("invalid parameter")
	}
	f.rootPath = path
	f.setRootPath(path)
	return nil
}

//////////////
//private func
//////////////

//get db obj by tag
func (f *Storage) getDbByTag(tag string) *DB {
	if tag == "" {
		return nil
	}
	v, ok := f.dbMap.Load(tag)
	if !ok || v == nil {
		return nil
	}
	obj, subOk := v.(*DB)
	if !subOk || obj == nil {
		return nil
	}
	return obj
}
