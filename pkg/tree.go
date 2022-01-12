package pkg

import (
	"fmt"
	"strings"
	"time"
)

type Node interface {
	Children() []Node
}

type CategoryNode struct {
	tag      int
	name     string
	children []Node
}

type ValueNode struct {
	name  string
	tag   string
	unit  *string
	value *float64
}

type TextNode struct {
	name string
	tag  string
	text string
}

type DurationNode struct {
	name     string
	tag      string
	duration time.Duration
}

func (c CategoryNode) Children() []Node {
	return c.children
}

func (c CategoryNode) String() string {
	return c.name
}

func (c *CategoryNode) FindCategory(tag int) *CategoryNode {
	if c.tag == tag {
		return c
	}
	for _, child := range c.children {
		catChild, ok := child.(*CategoryNode)
		if !ok {
			continue
		}
		r := catChild.FindCategory(tag)
		if r != nil {
			return r
		}
	}
	return nil
}

func (c *CategoryNode) AddChild(child Node) {
	c.children = append(c.children, child)
}

func (v ValueNode) Children() []Node {
	return []Node{}
}

func (v *ValueNode) String() string {
	strV := "n/a"
	if v.value != nil {
		strV = fmt.Sprintf("%f", *v.value)
	}

	if v.unit != nil {
		return fmt.Sprintf("%s: %s %s", v.name, strV, *v.unit)
	} else {
		return fmt.Sprintf("%s: %s", v.name, strV)
	}
}

func (t TextNode) Children() []Node {
	return []Node{}
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s: %s", t.name, t.text)
}

func (d DurationNode) Children() []Node {
	return []Node{}
}

func (d *DurationNode) String() string {
	return fmt.Sprintf("%s: %s", d.name, d.duration)
}

func PrintTree(root Node, level int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", level), root)
	for _, child := range root.Children() {
		PrintTree(child, level+1)
	}
}
