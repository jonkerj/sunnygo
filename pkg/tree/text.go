package tree

import "fmt"

type TextNode struct {
	tag  string
	name string
	text string
}

func NewTextNode(tag string, name string, text string) *TextNode {
	t := &TextNode{
		tag:  tag,
		name: name,
		text: text,
	}

	return t
}

func (t TextNode) Children() []Node {
	return []Node{}
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s: %s", t.name, t.text)
}
