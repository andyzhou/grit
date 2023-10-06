package testing

import "testing"

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

	//init db
	db, err := InitDB()
	if err != nil {
		b.Logf("init db failed, err:%v\n", err.Error())
		return
	}

	for n := 0; n < b.N; n++ {
		//get or set key
		rec, subErr := db.Get(key)
		b.Logf("get key:%v, err:%v\n", key, subErr)
		if subErr != nil || rec == nil {
			subErr = db.Set(key, val)
			b.Logf("set key:%v, err:%v\n", key, subErr)
		}else{
			b.Logf("get rec:%v\n", string(rec))
		}
	}

	for n := 0; n < b.N; n++ {
		bRet, _ := db.Exists(key)
		b.Logf("exists bRet:%v\n", bRet)
	}
}