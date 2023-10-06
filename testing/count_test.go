package testing

import "testing"

//count test
func TestCount(t *testing.T)  {
	var (
		key = "c1"
		hashKey = "hc1"
	)

	//init db
	db, err := InitDB()
	if err != nil {
		t.Logf("init db failed, err:%v\n", err.Error())
		return
	}

	//simple count
	for i := 0; i <= 1; i++ {
		v, subErr := db.IncBy(key, 1)
		t.Logf("inc v:%v, err:%v\n", v, subErr)
	}
	rec, _ := db.Get(key)
	t.Logf("get rec:%v\n", string(rec))

	//hash count
	err = db.HashIncBy(hashKey, map[string]int64{
		"up":2,
		"down":1,
	}, true)
	t.Logf("hash inc by, err:%v\n", err)

	recMap, err := db.GetHashCount(hashKey)
	t.Logf("hash count:%v, err:%v\n", recMap, err)
}

//count benchmark
func BenchmarkCount(b *testing.B) {
	var (
		key = "c1"
		hashKey = "hc1"
	)

	//reset timer
	b.ResetTimer()

	//init db
	db, err := InitDB()
	if err != nil {
		b.Logf("init db failed, err:%v\n", err.Error())
		return
	}

	//simple count
	for n := 0; n < b.N; n++ {
		//inc by count
		v, subErr := db.IncBy(key, 1)
		b.Logf("inc v:%v, err:%v\n", v, subErr)

		//get count
		rec, _ := db.Get(key)
		b.Logf("get rec:%v\n", string(rec))

		//hash count
		err = db.HashIncBy(hashKey, map[string]int64{
			"up":2,
			"down":1,
		}, true)
		b.Logf("hash inc by, err:%v\n", err)

		//get hash count
		recMap, subErrTwo := db.GetHashCount(hashKey)
		b.Logf("hash count:%v, err:%v\n", recMap, subErrTwo)
	}


}