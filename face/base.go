package face

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"sync"
)

/*
 * @author Andy Chow <diudiu8848@163.com>
 * base face
 * - all property and func should be private
 */

type Base struct {
	rootPath string
	db       *leveldb.DB //reference
	locker   sync.RWMutex
}

//del data
func (f *Base) remove(key []byte) error {
	//check
	if key == nil || len(key) <= 0 {
		return errors.New("invalid parameter")
	}
	if f.db == nil {
		return errors.New("db didn't init")
	}

	//delete data from db with locker
	f.locker.Lock()
	defer f.locker.Unlock()
	err := f.db.Delete(key, nil)
	return err
}

//read data
func (f *Base) read(key []byte) ([]byte, error) {
	//check
	if key == nil || len(key) <= 0 {
		return nil, errors.New("invalid parameter")
	}
	if f.db == nil {
		return nil, errors.New("db didn't init")
	}

	//get data from db with locker
	f.locker.Lock()
	defer f.locker.Unlock()
	data, err := f.db.Get(key, nil)
	return data, err
}

//check key are exists or not
func (f *Base) isExists(key []byte) (bool, error) {
	//check
	if key == nil {
		return false, errors.New("invalid parameter")
	}
	if f.db == nil {
		return false, errors.New("db didn't init")
	}

	//check from db with locker
	f.locker.Lock()
	defer f.locker.Unlock()
	bRet, err := f.db.Has(key, nil)
	return bRet, err
}
func (f *Base) mIsExists(keys ...[]byte) (map[interface{}]bool, error) {
	var (
		bRet bool
	)
	//check
	if keys == nil || len(keys) <= 0 {
		return nil, errors.New("invalid parameter")
	}
	if f.db == nil {
		return nil, errors.New("db didn't init")
	}

	//check from db with locker
	f.locker.Lock()
	defer f.locker.Unlock()
	result := make(map[interface{}]bool)
	for _, key := range keys {
		bRet, _ = f.db.Has(key, nil)
		result[key] = bRet
	}
	return result, nil
}

//save data
func (f *Base) save(key, data []byte) error {
	//check
	if key == nil || len(key) <= 0 ||
		data == nil || len(data) <= 0 {
		return errors.New("invalid parameter")
	}
	if f.db == nil {
		return errors.New("db didn't init")
	}

	//put data into db with locker
	f.locker.Lock()
	defer f.locker.Unlock()
	err := f.db.Put(key, data, nil)
	return err
}

//set db
func (f *Base) setDB(db *leveldb.DB) {
	f.locker.Lock()
	defer f.locker.Unlock()
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

	//detect and make dir
	_, err := os.Stat(dir)
	if err != nil {
		//dir not exist
		err = os.Mkdir(dir, 0755)
	}
	return err
}