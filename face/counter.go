package face

import (
	"encoding/json"
	"errors"
	"github.com/andyzhou/grit/data"
	"log"
)

/*
 * counter data face
 * - work on `DB` face
 * - atomic on process
 */

const (
	updateChanSize = 1024
)

//req data
type ReqOfUpdate struct {
	key string
	im map[string]int64
	resp chan error
}

//face info
type Counter struct {
	updateChan chan ReqOfUpdate
	closeChan chan struct{}
	started bool
	Base
}

//get
func (f *Counter) GetCount(key string) (map[string]int64, error) {
	//check
	if key == "" {
		return nil, errors.New("invalid parameter")
	}

	//get from db
	dataByte, err := f.Read([]byte(key))
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

//inc by
func (f *Counter) IncBy(key string, im map[string]int64) error {
	//check
	if key == "" || im == nil {
		return errors.New("invalid parameter")
	}

	//format request
	req := ReqOfUpdate {
		key: key,
		im: im,
		resp: make(chan error, 1),
	}

	//send to chan
	f.updateChan <- req

	//wait resp
	err, _ := <- req.resp
	return err
}

//close process
func (f *Counter) closeProcess() {
	close(f.closeChan)
}

//start update process
func (f *Counter) startCountProcess() {
	if f.started {
		return
	}
	//setup
	f.updateChan = make(chan ReqOfUpdate, updateChanSize)
	f.closeChan = make(chan struct{}, 1)

	//spawn process
	go f.runUpdateProcess()
	f.started = true
}

//update one data
func (f *Counter) updateOneData(req *ReqOfUpdate) error {
	var (
		cd data.CountData
	)
	//check
	if req == nil {
		return errors.New("invalid parameter")
	}
	//get record first
	dataByte, err := f.db.Get([]byte(req.key), nil)
	if err != nil {
		req.resp <- err
		return err
	}
	//check data
	if dataByte == nil {
		//init new
		cd = data.CountData{
			Fields: req.im,
		}
	}else{
		//decode
		cd = data.CountData{
			Fields: map[string]int64{},
		}
		err = json.Unmarshal(dataByte, &cd)
		if err != nil {
			req.resp <- err
			return err
		}
		//check and update
		for field, count := range req.im {
			v, ok := cd.Fields[field]
			if ok {
				v += count
			}else{
				v = count
			}
			if v <= 0 {
				v = 0
			}
			cd.Fields[field] = v
		}
	}
	//encode
	dataByte, err = json.Marshal(cd)
	if err != nil {
		req.resp <- err
		return err
	}
	//save into db
	err = f.Save([]byte(req.key), dataByte)
	req.resp <- err
	return err
}

//run update process
func (f *Counter) runUpdateProcess() {
	var (
		req ReqOfUpdate
		ok bool
	)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Counter:runUpdateProcess panic, err:%v\n", err)
		}
	}()

	for {
		select {
		case req, ok = <- f.updateChan:
			if ok {
				f.updateOneData(&req)
			}
		case <- f.closeChan:
			return
		}
	}
}