package parsers

import (
	"fmt"
	"path/filepath"
)

const (
	YAML_EXT       = ".yaml"
	YAML_EXT_SHORT = ".yml"
	JSON_EXT       = ".json"
)

func ParseByExtension(path string) (map[string]interface{}, error) {
	ext := filepath.Ext(path)
	switch ext {
	case JSON_EXT:
		return ParseJSON(path)
	case YAML_EXT, YAML_EXT_SHORT:
		return ParseYAML(path)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
