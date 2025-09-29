package models

type DiffNode struct {
	Key      string
	Status   string
	OldValue interface{}
	NewValue interface{}
	Children []*DiffNode
}
