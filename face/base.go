package face

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

/*
 * Copyright (c) 2023, AndyChow <diudiu8848@gmail.com>
 * All rights reserved.
 * base face
 * - all property and func should be private
 */

type Base struct {
	rootPath string
	db *leveldb.DB //reference
}

//del data
func (f *Base) remove(key []byte) error {
	//check
	if key == nil {
		return errors.New("invalid parameter")
	}
	if f.db == nil {
		return errors.New("db didn't init")
	}
	//try delete data from db
	err := f.db.Delete(key, nil)
	return err
}

//read data
func (f *Base) read(key []byte) ([]byte, error) {
	//check
	if key == nil {
		return nil, errors.New("invalid parameter")
	}
	if f.db == nil {
		return nil, errors.New("db didn't init")
	}
	//try get data from db
	data, err := f.db.Get(key, nil)
	return data, err
}

//save data
func (f *Base) save(key, data []byte) error {
	//check
	if key == nil || data == nil {
		return errors.New("invalid parameter")
	}
	if f.db == nil {
		return errors.New("db didn't init")
	}
	//try put data into db
	err := f.db.Put(key, data, nil)
	return err
}

//set db
func (f *Base) setDB(db *leveldb.DB) {
	f.db = db
}

//get root path
func (f *Base) getRootPath() string {
	return f.rootPath
}

//set root path
func (f *Base) setRootPath(path string) {
	f.rootPath = path
}

//check root dir
func (f *Base) checkDir(dir string) error {
	//check
	if dir == "" {
		return errors.New("invalid dir param")
	}
	//detect and try make dir
	_, err := os.Stat(dir)
	if err != nil {
		//dir not exist
		err = os.Mkdir(dir, 0755)
	}
	return err
}