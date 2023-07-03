package face

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andyzhou/grit/data"
	"github.com/andyzhou/grit/define"
	"log"
	"strconv"
	"sync"
)

/*
 * Copyright (c) 2023, AndyChow <diudiu8848@gmail.com>
 * All rights reserved.
 * counter data face
 * - work on `DB` face
 * - support update on queue
 */

//inter type
type (
	countReq struct {
		key string
		incVal int64 //second priority
		hashVal map[string]int64 //first priority
	}
)

//face info
type Counter struct {
	queueSize int
	hasInit bool
	queueRun bool
	locker sync.RWMutex
	reqChan chan countReq
	closeChan chan struct{}
	Base
}

//get hash count value
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

//get gen count value
func (f *Counter) GetCount(
		key string) (int64, error) {
	//check
	if key == "" {
		return 0, errors.New("invalid parameter")
	}

	//get from db
	dataByte, err := f.read([]byte(key))
	if err != nil {
		return 0, err
	}
	//decode
	v, err := strconv.ParseInt(string(dataByte), 10, 64)
	return v, err
}

//inc by one key multi fields count
func (f *Counter) HashIncBy(
		key string,
		im map[string]int64,
		useQueues ...bool) error {
	var (
		useQueue bool
		err error
	)
	//check
	if key == "" || im == nil {
		return errors.New("invalid parameter")
	}
	if useQueues != nil && len(useQueues) > 0 {
		useQueue = useQueues[0]
	}
	if useQueue {
		//run in queue
		err = f.sendToQueue(key, 0, im)
		return err
	}

	//call directly
	err = f.interHashIncBy(key, im)
	return err
}

//inc by for one key one count
//if use queue mod, new value will be 0
//return new value, error
func (f *Counter) IncBy(
		key string,
		val int64,
		useQueues ...bool) (int64, error) {
	var (
		useQueue bool
		newVal int64
		err error
	)
	//check
	if key == "" || val == 0 {
		return 0, errors.New("invalid parameter")
	}
	if useQueues != nil && len(useQueues) > 0 {
		useQueue = useQueues[0]
	}
	if useQueue {
		//run in queue
		err = f.sendToQueue(key, val, nil)
		return 0, err
	}
	//call directly
	newVal, err = f.interIncBy(key, val)
	return newVal, err
}

//set inter queue size
func (f *Counter) SetQueueSize(size int) {
	if size <= 0 {
		return
	}
	f.queueSize = size
}

//////////////
//private func
//////////////

//send to queue
func (f *Counter) sendToQueue(
		key string,
		val int64,
		hv map[string]int64) error {
	//check
	if key == "" || (val == 0 && (hv == nil || len(hv) <= 0)) {
		return errors.New("invalid parameter")
	}
	//check or init queue
	f.checkAndInitQueue()
	if !f.queueRun || f.reqChan == nil {
		return errors.New("inter queue hadn't run")
	}
	//init request
	req := countReq{
		key: key,
		incVal: val,
		hashVal: hv,
	}
	//send to chan async
	select {
	case f.reqChan <- req:
	}
	return nil
}

//inter hash count inc
func (f *Counter) interHashIncBy(
		key string,
		im map[string]int64) error {
	var (
		cd data.CountData
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
	//check and update
	for field, count := range im {
		v, ok := cd.Fields[field]
		if ok {
			v += count
		}else{
			v = count
		}
		//border check
		if v <= 0 {
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

//inter inc by for one key one count
//return new value, error
func (f *Counter) interIncBy(
		key string,
		val int64) (int64, error) {
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
	if countVal < 0 {
		countVal = 0
	}

	//save
	countValStr := fmt.Sprintf("%v", countVal)
	err := f.save([]byte(key), []byte(countValStr))
	if err != nil {
		return val, err
	}
	return countVal, nil
}

//update one counter
func (f *Counter) updateOneCounter(
		req *countReq) error {
	var (
		err error
	)
	//check
	if req == nil || req.key == "" {
		return errors.New("invalid parameter")
	}
	if req.hashVal != nil && len(req.hashVal) > 0 {
		//check hash count first
		err = f.interHashIncBy(req.key, req.hashVal)
	}else{
		//gen inc count
		_, err = f.interIncBy(req.key, req.incVal)
	}
	return err
}

//run inter worker
func (f *Counter) runWorker() {
	var (
		m any = nil
		req countReq
		isOk bool
	)
	//defer
	defer func() {
		if err := recover(); err != m {
			log.Printf("counter.runWorker panic, err:%v\n", err)
		}
		close(f.reqChan)
	}()

	//loop
	for {
		select {
		case req, isOk = <- f.reqChan:
			if isOk && &req != nil {
				//update one counter
			}
		case <- f.closeChan:
			return
		}
	}
}

//check and init queue
func (f *Counter) checkAndInitQueue() {
	if f.queueRun {
		return
	}
	if f.queueSize <= 0 {
		f.queueSize = define.CountQueueSize
	}
	//init inter worker
	f.reqChan = make(chan countReq, f.queueSize)
	f.closeChan = make(chan struct{}, 1)
	//spawn son process
	go f.runWorker()
	f.queueRun = true
}

//inter init
func (f *Counter) counterInit() {
	if f.hasInit {
		return
	}
	//setup
	f.locker = sync.RWMutex{}
	f.hasInit = true
}