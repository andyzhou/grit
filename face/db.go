package face

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

/*
 * db base face
 */

//face info
type DB struct {
	tag string
	db *leveldb.DB //db engine
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
	if f.db != nil {
		err := f.db.Close()
		return err
	}
	return nil
}

////get db
//func (f *DB) GetDB() (*leveldb.DB, error) {
//	if f.db == nil {
//		return nil, errors.New("db not opened")
//	}
//	return f.db, nil
//}

//open db
func (f *DB) OpenDB() error {
	//check
	if f.db != nil {
		return errors.New("db had opened")
	}

	//check dir
	fileDBPath := fmt.Sprintf("%v/%v", f.GetRootPath(), f.tag)
	err := f.CheckDir(fileDBPath)
	if err != nil {
		return err
	}

	//init level db
	db, err := leveldb.OpenFile(fileDBPath, nil)
	if err != nil {
		return err
	}
	f.db = db

	//set base db
	f.SetDB(db)
	return nil
}