package parsers

import (
	"encoding/json"
	"os"
)

func ParseJSON(path string) (map[string]interface{}, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw interface{}
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, err
	}

	switch v := raw.(type) {
	case map[string]interface{}:
		return v, nil
	case []interface{}:
		return map[string]interface{}{"root": v}, nil
	default:
		return map[string]interface{}{"root": v}, nil
	}
}
