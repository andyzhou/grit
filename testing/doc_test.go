package testing

import (
	"fmt"
	"github.com/andyzhou/grit/face"
	"log"
	"os"
	"testing"
)

//global variable
var (
	db *face.DB
	err error
)

//init
func init()  {
	//init db
	db, err = InitDB()
	if err != nil {
		log.Printf("init db failed, err:%v\n", err.Error())
		os.Exit(1)
		return
	}
}

//doc test
func TestDoc(t *testing.T) {
	var (
		key = "t1"
		val = "hello"
	)

	//init db
	db, err = InitDB()
	if err != nil {
		t.Logf("init db failed, err:%v\n", err.Error())
		return
	}

	//get or set key
	rec, subErr := db.Get(key)
	if subErr != nil || rec == nil {
		subErr = db.Set(key, val)
		t.Logf("set key:%v, err:%v\n", key, subErr)
	}else{
		t.Logf("get rec:%v\n", string(rec))
	}

	//check key is exists
	bRet, _ := db.Exists(key)
	t.Logf("exists bRet:%v\n", bRet)

	//close db
	if db != nil {
		db.CloseDB()
	}
}

//doc benchmark
func BenchmarkDoc(b *testing.B) {
	var (
		key string
		val string
		rec []byte
	)

	//reset timer
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		//get or set key
		key = fmt.Sprintf("test-%v", n)
		val = fmt.Sprintf("val-%v", n)

		//get key
		rec, err = db.Get(key)
		//b.Logf("get key:%v, rec:%v, err:%v\n", key, rec, subErr)
		if err != nil || rec == nil {
			b.Logf("get key %v failed, err:%v\n", key, err)

			//set new key
			err = db.Set(key, val)
			b.Logf("set key %v, err:%v\n", key, err)
		}
	}

	//for n := 0; n < b.N; n++ {
	//	db.Exists(key)
	//	//b.Logf("exists bRet:%v\n", bRet)
	//}

	//close db
	if db != nil {
		db.CloseDB()
	}
}