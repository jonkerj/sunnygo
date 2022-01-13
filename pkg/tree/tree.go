package tree

import (
	"fmt"
	"strings"
)

type Node interface {
	Children() []Node
}

func PrintTree(root Node, level int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", level), root)
	for _, child := range root.Children() {
		PrintTree(child, level+1)
	}
}
