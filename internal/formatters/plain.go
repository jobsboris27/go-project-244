package formatters

import (
	models "code/internal/models"
	"fmt"
	"strings"
)

func RenderPlain(diffNodes []*models.DiffNode, path string) string {
	var result strings.Builder

	if len(diffNodes) == 0 {
		return ""
	}

	for _, node := range diffNodes {
		currentPath := buildPath(path, node.Key)

		switch node.Status {
		case "added":
			result.WriteString(fmt.Sprintf("Property '%s' was added with value: %s\n",
				currentPath, formatPlainValue(node.NewValue)))

		case "removed":
			result.WriteString(fmt.Sprintf("Property '%s' was removed\n", currentPath))

		case "modified":
			result.WriteString(fmt.Sprintf("Property '%s' was updated. From %s to %s\n",
				currentPath, formatPlainValue(node.OldValue), formatPlainValue(node.NewValue)))

		case "nested":
			nestedResult := RenderPlain(node.Children, currentPath)
			if nestedResult != "" {
				result.WriteString(nestedResult)
				result.WriteString("\n")
			}
		}
	}

	return strings.TrimSpace(result.String())
}

func buildPath(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

func formatPlainValue(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return "''"
		}
		return fmt.Sprintf("'%s'", v)
	case map[string]interface{}, map[interface{}]interface{}, []interface{}:
		return "[complex value]"
	default:
		return fmt.Sprintf("%v", v)
	}
}
