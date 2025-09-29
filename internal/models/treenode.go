package models

type TreeNode struct {
	Key      string
	Value    interface{}
	Children []*TreeNode
}
