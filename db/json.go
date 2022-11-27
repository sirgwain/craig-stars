package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

// helper to convert an item into JSON
func valueJSON(item interface{}) (driver.Value, error) {
	if isNil(item) {
		return nil, nil
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// helper to scan a text JSON column back into a struct
func scanJSON(src interface{}, dest interface{}) error {
	if src == nil {
		// leave empty
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, dest)
	case string:
		return json.Unmarshal([]byte(v), dest)
	}
	return errors.New("type assertion failed")
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
