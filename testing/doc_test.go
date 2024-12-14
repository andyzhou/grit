package testing

import (
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
	db, err := InitDB()
	if err != nil {
		t.Logf("init db failed, err:%v\n", err.Error())
		return
	}

	//get or set key
	rec, err := db.Get(key)
	if err != nil || rec == nil {
		err = db.Set(key, val)
		t.Logf("set key:%v, err:%v\n", key, err)
	}else{
		t.Logf("get rec:%v\n", string(rec))
	}

	//check key is exists
	bRet, _ := db.Exists(key)
	t.Logf("exists bRet:%v\n", bRet)
}

//doc benchmark
func BenchmarkDoc(b *testing.B) {
	var (
		key = "t1"
		val = "hello"
	)

	//reset timer
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		//get or set key
		rec, subErr := db.Get(key)
		//b.Logf("get key:%v, rec:%v, err:%v\n", key, rec, subErr)
		if subErr != nil || rec == nil {
			subErr = db.Set(key, val)
			if subErr != nil {
				b.Logf("set key:%v, err:%v\n", key, subErr)
			}
		}else{
			//b.Logf("get rec:%v\n", string(rec))
		}
	}

	//for n := 0; n < b.N; n++ {
	//	db.Exists(key)
	//	//b.Logf("exists bRet:%v\n", bRet)
	//}
}