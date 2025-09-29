package formatters

import (
	"encoding/json"
)

func RenderJSON(diffNodes []*DiffNode) string {
	jsonData := convertToJSONFormat(diffNodes)
	result, _ := json.MarshalIndent(jsonData, "", "  ")
	return string(result)
}

func convertToJSONFormat(diffNodes []*DiffNode) []map[string]interface{} {
	var result []map[string]interface{}

	for _, node := range diffNodes {
		switch node.Status {
		case "added":
			result = append(result, map[string]interface{}{
				"key":    node.Key,
				"type":   "added",
				"value":  node.NewValue,
			})

		case "removed":
			result = append(result, map[string]interface{}{
				"key":   node.Key,
				"type":  "removed",
				"value": node.OldValue,
			})

		case "unchanged":
			result = append(result, map[string]interface{}{
				"key":   node.Key,
				"type":  "unchanged",
				"value": node.OldValue,
			})

		case "modified":
			result = append(result, map[string]interface{}{
				"key":      node.Key,
				"type":     "updated",
				"oldValue": node.OldValue,
				"newValue": node.NewValue,
			})

		case "nested":
			result = append(result, map[string]interface{}{
				"key":      node.Key,
				"type":     "nested",
				"children": convertToJSONFormat(node.Children),
			})
		}
	}

	return result
}