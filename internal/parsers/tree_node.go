package parsers

import (
	"code/internal/models"
	"fmt"
	"sort"
)

func СonvertMapToTree(data map[string]interface{}) *models.TreeNode {
	root := &models.TreeNode{Key: "root", Children: []*models.TreeNode{}}

	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			childTree := СonvertMapToTree(v)
			childNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case map[interface{}]interface{}:
			convertedMap := make(map[string]interface{})
			for k, val := range v {
				if strKey, ok := k.(string); ok {
					convertedMap[strKey] = val
				}
			}
			childTree := СonvertMapToTree(convertedMap)
			childNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case []interface{}:
			arrayNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: []*models.TreeNode{},
			}

			for i, item := range v {
				itemNode := &models.TreeNode{
					Key:   fmt.Sprintf("[%d]", i),
					Value: item,
				}

				switch v := item.(type) {
				case map[string]interface{}:
					itemNode = СonvertMapToTree(v)
					itemNode.Key = fmt.Sprintf("[%d]", i)
				case []interface{}:
					itemNode = СonvertMapToTree(map[string]interface{}{
						fmt.Sprintf("[%d]", i): v,
					})
				default:
					itemNode = &models.TreeNode{
						Key:   fmt.Sprintf("[%d]", i),
						Value: item,
					}
				}

				arrayNode.Children = append(arrayNode.Children, itemNode)
			}
			root.Children = append(root.Children, arrayNode)
		default:
			leafNode := &models.TreeNode{
				Key:   key,
				Value: v,
			}
			root.Children = append(root.Children, leafNode)
		}
	}

	return root
}

func GetDiff(tree1, tree2 *models.TreeNode) []*models.DiffNode {
	var diff []*models.DiffNode

	allKeys := collectAllKeys(tree1, tree2)
	sort.Strings(allKeys)

	for _, key := range allKeys {
		node1 := findChildByKey(tree1, key)
		node2 := findChildByKey(tree2, key)

		diffNode := &models.DiffNode{Key: key}

		switch {
		case node1 == nil && node2 != nil:
			diffNode.Status = "added"
			if hasChildren(node2) {
				diffNode.NewValue = reconstructObject(node2)
			} else {
				diffNode.NewValue = node2.Value
			}

		case node1 != nil && node2 == nil:
			diffNode.Status = "removed"
			if hasChildren(node1) {
				diffNode.OldValue = reconstructObject(node1)
			} else {
				diffNode.OldValue = node1.Value
			}

		case node1 != nil && node2 != nil:
			if hasChildren(node1) && hasChildren(node2) {
				diffNode.Status = "nested"
				diffNode.Children = GetDiff(node1, node2)
			} else if !hasChildren(node1) && !hasChildren(node2) && areValuesEqual(node1, node2) {
				diffNode.Status = "unchanged"
				diffNode.OldValue = node1.Value
			} else {
				diffNode.Status = "modified"
				if hasChildren(node1) {
					diffNode.OldValue = reconstructObject(node1)
				} else {
					diffNode.OldValue = node1.Value
				}
				if hasChildren(node2) {
					diffNode.NewValue = reconstructObject(node2)
				} else {
					diffNode.NewValue = node2.Value
				}
			}
		}

		diff = append(diff, diffNode)
	}

	return diff
}

func collectAllKeys(tree1, tree2 *models.TreeNode) []string {
	keys := make(map[string]bool)

	if tree1 != nil {
		for _, child := range tree1.Children {
			keys[child.Key] = true
		}
	}
	if tree2 != nil {
		for _, child := range tree2.Children {
			keys[child.Key] = true
		}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}

func findChildByKey(tree *models.TreeNode, key string) *models.TreeNode {
	if tree == nil {
		return nil
	}
	for _, child := range tree.Children {
		if child.Key == key {
			return child
		}
	}
	return nil
}

func areValuesEqual(node1, node2 *models.TreeNode) bool {
	if node1 == nil || node2 == nil {
		return false
	}
	return fmt.Sprintf("%v", node1.Value) == fmt.Sprintf("%v", node2.Value)
}

func hasChildren(node *models.TreeNode) bool {
	return node != nil && len(node.Children) > 0
}

func reconstructObject(node *models.TreeNode) interface{} {
	if !hasChildren(node) {
		return node.Value
	}

	result := make(map[string]interface{})
	for _, child := range node.Children {
		result[child.Key] = reconstructObject(child)
	}
	return result
}
