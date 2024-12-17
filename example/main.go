package main

import (
	"fmt"
	"github.com/andyzhou/grit"
	"github.com/andyzhou/grit/face"
	"os"
	"sync"
	"time"
)

/*
 * @author Andy Chow <diudiu8848@163.com>
 * example code
 */

//test doc
func testDoc(db *face.DB) {
	var (
		key = "t1"
		val = "hello"
	)

	//get or set key
	rec, err := db.Get(key)
	if err != nil || rec == nil {
		err = db.Set(key, val)
		fmt.Println("set err:", err)
	}else{
		fmt.Println("get rec:", string(rec))
	}

	//check key is exists
	bRet, err := db.Exists(key)
	fmt.Println("bRet:", bRet, ", err:", err)
}

//test counter
func testCounter(db *face.DB) {
	var (
		key = "c1"
		hashKey = "hc1"
	)

	//simple count
	for i := 0; i <= 1; i++ {
		v, err := db.IncBy(key, 1)
		fmt.Println("v:", v, ", err:", err)
	}
	rec, _ := db.Get(key)
	fmt.Println("rec:", string(rec))

	//hash count
	err := db.HashIncBy(hashKey, map[string]int64{
		"up":2,
		"down":1,
	}, true)
	fmt.Println("hash inc by, err:", err)

	recMap, err := db.GetHashCount(hashKey)
	fmt.Println("hash count:", recMap, ", err:", err)
}

//force done
func forceDone(num int, wg *sync.WaitGroup)  {
 for i := 0; i < num; i++ {
	 wg.Done()
 }
}

func main() {
	var (
		wg sync.WaitGroup
		dbTag = "test"
	)

	//try open db
	g := grit.GetGrit()
	db, err := g.OpenDB(dbTag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//defer close db
	defer db.CloseDB()

	fmt.Println("start..")
	wg.Add(2)

	//test
	go testDoc(db)
	go testCounter(db)

	//get dbs
	dbs := g.GetDBs()
	fmt.Println("dbs:",  dbs)

	//force end
	sf := func() {
		forceDone(2, &wg)
	}
	time.AfterFunc(time.Second * 2, sf)

	wg.Wait()
	fmt.Println("done!")
}
