package face

import (
	"errors"
	"github.com/andyzhou/grit/define"
	"runtime"
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
	dbMap    map[string]*DB //tag -> *DB
	Base
	sync.RWMutex
}

//construct
func NewStorage() *Storage {
	//self init
	this := &Storage{
		rootPath: define.DefaultDBRootPath,
		dbMap: map[string]*DB{},
	}
	return this
}

//quit
func (f *Storage) Quit() {
	//close all db
	f.Lock()
	defer f.Unlock()
	for k, v := range f.dbMap {
		v.CloseDB()
		delete(f.dbMap, k)
	}

	//gc opt
	runtime.GC()
}

//get dbs
func (f *Storage) GetDBs() []string {
	f.Lock()
	defer f.Unlock()
	result := make([]string, 0)
	for tag, _ := range f.dbMap {
		if tag == "" {
			continue
		}
		result = append(result, tag)
	}
	return result
}

//close db
func (f *Storage) CloseDB(tag string) error {
	//check
	if tag == "" {
		return errors.New("invalid parameter")
	}

	//hit cache first
	db, _ := f.getDbByTag(tag)
	if db == nil {
		return errors.New("no db obj by tag")
	}

	//try close db
	err := db.CloseDB()
	if err != nil {
		return err
	}

	//remove from map with locker
	f.Lock()
	defer f.Unlock()
	delete(f.dbMap, tag)

	//gc opt
	if len(f.dbMap) <= 0 {
		runtime.GC()
	}
	return nil
}

//get db
func (f *Storage) GetDB(tag string) (*DB, error) {
	//check
	if tag == "" {
		return nil, errors.New("invalid parameter")
	}

	//hit cache first
	db, err := f.getDbByTag(tag)
	if err != nil {
		return nil, err
	}
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
	db, err := f.getDbByTag(tag)
	if err != nil || db != nil {
		return db, nil
	}

	//try open db
	db = NewDB(tag)
	db.setRootPath(f.rootPath)
	err = db.OpenDB()
	if err != nil {
		return nil, err
	}

	//sync into map with locker
	f.Lock()
	defer f.Unlock()
	f.dbMap[tag] = db
	return db, nil
}

//set db path
func (f *Storage) SetDBPath(path string) error {
	//check
	if path == "" {
		return errors.New("invalid parameter")
	}

	//set root path
	f.rootPath = path
	f.setRootPath(path)
	return nil
}

//////////////
//private func
//////////////

//get db obj by tag
func (f *Storage) getDbByTag(tag string) (*DB, error){
	//check
	if tag == "" {
		return nil, errors.New("invalid parameter")
	}

	//get with locker
	f.Lock()
	defer f.Unlock()
	v, ok := f.dbMap[tag]
	if !ok || v == nil {
		return nil, nil
	}
	return v, nil
}
