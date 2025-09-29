package formatters

import (
	"fmt"
	"sort"
	"strings"
)

func RenderStylish(diffNodes []*DiffNode, depth int) string {
	var result strings.Builder

	if depth == 0 {
		result.WriteString("{\n")
	}

	for _, node := range diffNodes {
		switch node.Status {
		case "unchanged":
			indent := strings.Repeat(" ", depth*4+4)
			result.WriteString(fmt.Sprintf("%s%s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))

		case "added":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s+ %s: %s\n", indent, node.Key, formatValue(node.NewValue, depth+1)))

		case "removed":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))

		case "modified":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))
			result.WriteString(fmt.Sprintf("%s+ %s: %s\n", indent, node.Key, formatValue(node.NewValue, depth+1)))

		case "nested":
			indent := strings.Repeat(" ", depth*4+4)
			result.WriteString(fmt.Sprintf("%s%s: {\n", indent, node.Key))
			result.WriteString(RenderStylish(node.Children, depth+1))
			result.WriteString(fmt.Sprintf("%s}\n", indent))
		}
	}

	if depth == 0 {
		result.WriteString("}")
	}

	return result.String()
}

func formatValue(value interface{}, depth int) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return ""
		}
		return v
	case map[string]interface{}:
		return formatObject(v, depth)
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, val := range v {
			if strKey, ok := k.(string); ok {
				converted[strKey] = val
			}
		}
		return formatObject(converted, depth)
	case []interface{}:
		return formatArray(v, depth)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatObject(obj map[string]interface{}, depth int) string {
	if len(obj) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		indent := strings.Repeat(" ", depth*4+4)
		result.WriteString(fmt.Sprintf("%s%s: %s\n", indent, key, formatValue(obj[key], depth+1)))
	}

	indent := strings.Repeat(" ", depth*4)
	result.WriteString(fmt.Sprintf("%s}", indent))
	return result.String()
}

func formatArray(arr []interface{}, depth int) string {
	if len(arr) == 0 {
		return "[]"
	}

	var result strings.Builder
	result.WriteString("[\n")

	for i, item := range arr {
		indent := strings.Repeat(" ", depth*4+4)
		result.WriteString(fmt.Sprintf("%s%s\n", indent, formatValue(item, depth+1)))
		if i < len(arr)-1 {
			result.WriteString(",")
		}
	}

	indent := strings.Repeat(" ", depth*4)
	result.WriteString(fmt.Sprintf("%s]", indent))
	return result.String()
}