package util

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

func Json2Map(jsonStr string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, errors.Wrap(err, "util:json:")
	}
	return m, nil
}

func Map2Json(m map[string]interface{}) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return "", errors.Wrap(err, "util:json:")
	}
	return string(jsonByte), nil
}

// C style: return Pstructer
func Json2Struct(Pstructer interface{}, jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), Pstructer)
	if err != nil {
		return errors.Wrap(err, "util:json:")
	}
	return nil
}

func Struct2Json(structer interface{}) (string, error) {
	jsonBytes, err := json.Marshal(structer)
	if err != nil {
		return "", errors.Wrap(err, "util:json:")
	}
	return string(jsonBytes), err
}

func Map2Struct(Pstructer interface{}, m map[string]interface{}) error {
	s, err := Map2Json(m)
	if err != nil {
		return errors.Wrap(err, "util:json:")
	}

	Json2Struct(Pstructer, s)
	return nil
}

func Struct2Map(structer interface{}) (map[string]interface{}, error) {
	s, err := Struct2Json(structer)
	if err != nil {
		return nil, errors.Wrap(err, "util:json:")
	}

	m, err := Json2Map(s)
	if err != nil {
		return nil, errors.Wrap(err, "util:json:")
	}
	return m, nil
}

// the json comment in struct is ignored
// NOT for PROD!!!
func Struct2MapV2(structer interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	t := reflect.TypeOf(structer)
	v := reflect.ValueOf(structer)

	fieldNum := t.NumField()

	for i := 0; i < fieldNum; i++ {
		m[t.Field(i).Name] = v.Field(i).Interface()
	}
	return m, nil
}
