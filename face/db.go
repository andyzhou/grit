package face

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

/*
 * @author Andy Chow <diudiu8848@163.com>
 * db base face
 */

//face info
type DB struct {
	tag string
	db  *leveldb.DB //db engine
	Doc
	Counter
	Base
}

//construct
func NewDB(tag string) *DB {
	this := &DB{
		tag: tag,
	}
	return this
}

//close db
func (f *DB) CloseDB() error {
	var (
		err error
	)
	if f.db != nil {
		err = f.db.Close()
	}
	return err
}

//open db
func (f *DB) OpenDB() error {
	//check
	if f.db != nil {
		return errors.New("db had opened")
	}

	//check dir
	fileDBPath := fmt.Sprintf("%v/%v", f.getRootPath(), f.tag)
	err := f.checkDir(fileDBPath)
	if err != nil {
		return err
	}

	//init level db
	db, subErr := leveldb.OpenFile(fileDBPath, nil)
	if subErr != nil {
		return subErr
	}
	f.db = db

	//set base db
	f.setDB(db)
	f.Doc.setDB(db)
	f.Counter.setDB(db)

	//counter init
	f.counterInit()
	return nil
}