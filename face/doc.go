package face

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/andyzhou/grit/define"
)

/*
 * doc data face
 * - work on `DB` face
 * - multi read and write
 */

//face info
type Doc struct {
	Base
}

//del
func (f *Doc) MDel(keys ...string) error {
	var (
		err error
	)
	if keys == nil || len(keys) <= 0 {
		return errors.New("invalid parameter")
	}
	for _, key := range keys {
		err = f.Del(key)
		if err != nil {
			break
		}
	}
	return err
}

func (f *Doc) Del(key string) error {
	if key == "" {
		return errors.New("invalid parameter")
	}
	err := f.remove([]byte(key))
	return err
}

//get
func (f *Doc) MGet(keys ...string) (map[string][]byte, error) {
	if keys == nil || len(keys) <= 0 {
		return nil, errors.New("invalid parameter")
	}
	result := map[string][]byte{}
	for _, key := range keys {
		data, err := f.Get(key)
		if err != nil || data == nil {
			continue
		}
		result[key] = data
	}
	return result, nil
}

func (f *Doc) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("invalid parameter")
	}
	data, err := f.read([]byte(key))
	if err != nil &&
       err.Error() == define.ErrOfNotFound {
		return nil, nil
	}
	return data, err
}

//set
func (f *Doc) Set(key string, value interface{}) error {
	var (
		bf bytes.Buffer
	)
	if key == "" || value == nil {
		return errors.New("invalid parameter")
	}
	defer bf.Reset()
	if v, ok := value.([]byte); ok {
		bf.Write(v)
	}else{
		bf.WriteString(fmt.Sprintf("%v", value))
	}
	err := f.save([]byte(key), bf.Bytes())
	return err
}

//exists
func (f *Doc) Exists(key string) (bool, error) {
	if key == "" {
		return false, errors.New("invalid parameter")
	}
	bRet, err := f.db.Has([]byte(key), nil)
	return bRet, err
}
