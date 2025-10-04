package formatters

import (
	models "code/internal/models"
	"encoding/json"
	"fmt"
)

func RenderJSON(diffNodes []*models.DiffNode) string {
	jsonData := convertToJSONFormat(diffNodes)

	resultMap := map[string]interface{}{
		"diff": jsonData,
	}

	result, err := json.MarshalIndent(resultMap, "", "  ")

	if err != nil {
		fmt.Println("render json %w", err)
		return ""
	}
	return string(result)
}

func convertToJSONFormat(diffNodes []*models.DiffNode) []map[string]interface{} {
	var result []map[string]interface{}

	for _, node := range diffNodes {
		switch node.Status {
		case ADDED:
			result = append(result, map[string]interface{}{
				"key":   node.Key,
				"type":  ADDED,
				"value": node.NewValue,
			})

		case REMOVED:
			result = append(result, map[string]interface{}{
				"key":   node.Key,
				"type":  REMOVED,
				"value": node.OldValue,
			})

		case UNCHANGED:
			result = append(result, map[string]interface{}{
				"key":   node.Key,
				"type":  UNCHANGED,
				"value": node.OldValue,
			})

		case MODIFIED:
			result = append(result, map[string]interface{}{
				"key":      node.Key,
				"type":     UPDATED,
				"oldValue": node.OldValue,
				"newValue": node.NewValue,
			})

		case NESTED:
			result = append(result, map[string]interface{}{
				"key":      node.Key,
				"type":     NESTED,
				"children": convertToJSONFormat(node.Children),
			})
		}
	}

	return result
}
