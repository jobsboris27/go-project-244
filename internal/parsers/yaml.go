package parsers

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

func ParseYAML(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var raw interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	switch v := raw.(type) {
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, val := range v {
			if strKey, ok := k.(string); ok {
				converted[strKey] = val
			} else {
				converted[fmt.Sprintf("%v", k)] = val
			}
		}
		return converted, nil
	case map[string]interface{}:
		return v, nil
	case []interface{}:
		return map[string]interface{}{"root": v}, nil
	default:
		return map[string]interface{}{"root": v}, nil
	}
}
