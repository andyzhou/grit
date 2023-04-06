package face

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

/*
 * base face
 */

type Base struct {
	RootPath string
	db *leveldb.DB //reference
}

//del data
func (f *Base) Remove(key []byte) error {
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
func (f *Base) Read(key []byte) ([]byte, error) {
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
func (f *Base) Save(key, data []byte) error {
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
func (f *Base) SetDB(db *leveldb.DB) {
	f.db = db
}

//get root path
func (f *Base) GetRootPath() string {
	return f.RootPath
}

//set root path
func (f *Base) SetRootPath(path string) {
	f.RootPath = path
}

//check root dir
func (f *Base) CheckDir(dir string) error {
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