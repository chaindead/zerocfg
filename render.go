package zfg

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

// Converts dotted keys into nested maps and allows comment insertion.
func buildYAMLNode(input map[string]string, comments map[string]string) (*ast.MappingNode, error) {
	// Create root mapping node
	root := &ast.MappingNode{}

	// Build nested structure for each key
	for key, value := range input {
		// Split the key by dots
		parts := strings.Split(key, ".")

		// Start at root node
		current := root

		// Create nested maps for each part except last
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]

			// Try to find existing mapping node
			var nextNode *ast.MappingNode
			found := false

			for _, mapValue := range current.Values {
				if keyNode, ok := mapValue.Key.(*ast.StringNode); ok {
					if keyNode.Value == part {
						if mapNode, ok := mapValue.Value.(*ast.MappingNode); ok {
							nextNode = mapNode
							found = true
							break
						}
					}
				}
			}

			// Create new mapping node if not found
			if !found {
				keyNode := &ast.StringNode{Value: part}
				nextNode = &ast.MappingNode{}
				current.Values = append(current.Values, &ast.MappingValueNode{
					Key:   keyNode,
					Value: nextNode,
				})
			}

			current = nextNode
		}

		// Add the final key-value pair
		lastKey := parts[len(parts)-1]
		keyNode := &ast.StringNode{Value: lastKey}
		valueNode := &ast.StringNode{Value: value}

		// Create mapping value node
		mvNode := &ast.MappingValueNode{
			Key:   keyNode,
			Value: valueNode,
		}

		// Add comment if exists
		if comment, exists := comments[key]; exists {
			commentToken := &token.Token{
				Type:  token.CommentType,
				Value: "# " + comment,
			}
			commentGroup := &ast.CommentGroupNode{
				Comments: []*ast.CommentNode{
					{Token: commentToken},
				},
			}
			mvNode.Comment = commentGroup
		}

		current.Values = append(current.Values, mvNode)
	}

	return root, nil
}
