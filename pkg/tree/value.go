package tree

import "fmt"

type ValueNode struct {
	tag   string
	name  string
	unit  *string
	value *float64
}

func NewValueNode(tag string, name string) *ValueNode {
	v := &ValueNode{
		tag:  tag,
		name: name,
	}
	return v
}

func (v ValueNode) Children() []Node {
	return []Node{}
}

func (v *ValueNode) SetUnit(unit *string) {
	v.unit = unit
}

func (v *ValueNode) SetValue(value float64) {
	v.value = &value
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
