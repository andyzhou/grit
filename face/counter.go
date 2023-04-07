package face

import (
	"encoding/json"
	"errors"
	"fmt"
	"grit/data"
	"strconv"
	"sync"
)

/*
 * counter data face
 * - work on `DB` face
 * - atomic on process
 */

//face info
type Counter struct {
	hasInit bool
	locker sync.RWMutex
	Base
}

//get
func (f *Counter) GetHashCount(
			key string) (map[string]int64, error) {
	//check
	if key == "" {
		return nil, errors.New("invalid parameter")
	}

	//get from db
	dataByte, err := f.read([]byte(key))
	if err != nil {
		return nil, err
	}
	//decode
	cd := data.CountData{
		Fields: map[string]int64{},
	}
	err = json.Unmarshal(dataByte, &cd)
	if err != nil {
		return nil, err
	}
	return cd.Fields, nil
}

//inc by one key multi fields count
func (f *Counter) HashIncBy(
			key string,
			im map[string]int64,
			borderChecks ...bool) error {
	var (
		cd data.CountData
		borderCheck bool
	)
	//check
	if key == "" || im == nil {
		return errors.New("invalid parameter")
	}

	//check and update hash count with locker
	f.locker.Lock()
	defer f.locker.Unlock()

	//get record first
	dataByte, _ := f.db.Get([]byte(key), nil)
	//check data
	if dataByte == nil {
		//init new
		cd = data.CountData{
			Fields: im,
		}
	}else {
		//decode
		cd = data.CountData{
			Fields: map[string]int64{},
		}
		err := json.Unmarshal(dataByte, &cd)
		if err != nil {
			return err
		}
	}

	//border value check
	if borderChecks != nil && len(borderChecks) > 0 {
		borderCheck = borderChecks[0]
	}
	//check and update
	for field, count := range im {
		v, ok := cd.Fields[field]
		if ok {
			v += count
		}else{
			v = count
		}
		//border check
		if borderCheck && v <= 0 {
			v = 0
		}
		cd.Fields[field] = v
	}
	//encode
	dataByte, err := json.Marshal(cd)
	if err != nil {
		return err
	}
	//save into db
	err = f.save([]byte(key), dataByte)
	return err
}

//inc by for one key one count
//return new value, error
func (f *Counter) IncBy(
			key string,
			val int64,
			borderChecks ...bool) (int64, error) {
	var (
		countVal int64
	)
	//check
	if key == "" {
		return val, errors.New("invalid parameter")
	}
	//get and update with locker
	f.locker.Lock()
	defer f.locker.Unlock()

	//get first
	dataByte, _ := f.read([]byte(key))
	if dataByte != nil {
		v, err := strconv.ParseInt(string(dataByte), 10, 64)
		if err != nil {
			return val, errors.New("invalid key format")
		}
		countVal = v + val
	}else{
		countVal = val
	}

	//border value check
	if borderChecks != nil && len(borderChecks) > 0 {
		if borderChecks[0] && countVal < 0 {
			countVal = 0
		}
	}

	//save
	countValStr := fmt.Sprintf("%v", countVal)
	err := f.save([]byte(key), []byte(countValStr))
	if err != nil {
		return val, err
	}
	return countVal, nil
}

//////////////
//private func
//////////////

//inter init
func (f *Counter) counterInit() {
	if f.hasInit {
		return
	}
	//setup
	f.locker = sync.RWMutex{}

	//spawn process
	//go f.runUpdateProcess()
	f.hasInit = true
}