package generated

import (
	"database/sql/driver"
	"errors"

	json "github.com/json-iterator/go"
)

// helper to convert an item into JSON
func valueJSON(item interface{}) (driver.Value, error) {
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
