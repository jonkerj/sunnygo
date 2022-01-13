package tree

type CategoryNode struct {
	tag      int
	name     string
	children []Node
}

func NewCategoryNode(tag int, name string) *CategoryNode {
	newCategory := &CategoryNode{
		tag:      tag,
		name:     name,
		children: []Node{},
	}

	return newCategory
}

func (c *CategoryNode) AddChild(child Node) {
	c.children = append(c.children, child)
}

func (c CategoryNode) Children() []Node {
	return c.children
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

func (c CategoryNode) String() string {
	return c.name
}
