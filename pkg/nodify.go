package pkg

import (
	"encoding/json"
	"fmt"
	"time"
)

func (m *Meta) EnsureHierarchy(hierarchy []int, root *CategoryNode) error {
	parent := root
	for _, categoryTag := range hierarchy {
		r := parent.FindCategory(categoryTag)
		if r != nil {
			continue
		}

		// this category is not yet present in the tree. Create the node
		translation, err := m.GetTranslation(categoryTag)
		if err != nil {
			return fmt.Errorf("could not complete category tree due to missing translation: %v", err)
		}
		newCategory := &CategoryNode{
			tag:      categoryTag,
			name:     *translation,
			children: []Node{},
		}
		parent.children = append(parent.children, newCategory)
		parent = newCategory
	}

	return nil
}

func (i *IntValue) Nodify(tag string, meta *Meta) (Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	v := &ValueNode{
		tag:   tag,
		name:  *name,
		unit:  nil,
		value: nil,
	}

	if model.Unit != nil {
		unit, err := meta.GetTranslation(*model.Unit)
		if err != nil {
			return nil, fmt.Errorf("could not translation for unit %d: %w", *model.Unit, err)
		}
		v.unit = unit
	}

	val := float64(i.Val) * (*model.Scale)
	v.value = &val

	return v, nil
}

func (s *StringValue) Nodify(tag string, meta *Meta) (Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	t := &TextNode{
		tag:  tag,
		name: *name,
		text: s.Val,
	}

	return t, nil
}

func (d *DurationValue) Nodify(tag string, meta *Meta) (Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	n := &DurationNode{
		tag:      tag,
		name:     *name,
		duration: time.Duration(d.Val) * time.Second,
	}

	return n, nil
}

func (t *TagListValue) Nodify(tag string, meta *Meta) (Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	if len(t.Val) == 0 {
		return nil, fmt.Errorf("there are no tags in this result")
	}

	text, err := meta.GetTranslation(t.Val[0].Tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	n := &TextNode{
		tag:  tag,
		name: *name,
		text: *text,
	}

	return n, nil
}

func (m *Meta) NodifyAllValues(deviceId string, response *ResultReponse) (*CategoryNode, error) {
	fields, ok := response.Result[deviceId]
	if !ok {
		return nil, fmt.Errorf("device %s not found in values", deviceId)
	}

	root := &CategoryNode{
		name:     "root",
		children: []Node{},
	}

	for tag, field := range fields {
		model, err := m.GetModel(tag)
		if err != nil {
			return nil, fmt.Errorf("error getting model for tag %s: %w", tag, err)
		}

		var intf Nodifyable = nil

		switch model.DataFormat {
		case 0, 1, 2, 3, 26:
			intf = &IntValue{}
		case 7:
			intf = &DurationValue{}
		case 8:
			intf = &StringValue{}
		case 18:
			intf = &TagListValue{}
		default:
			fmt.Printf("skipping field %s (format %d) for now\n", tag, model.DataFormat)
			continue
		}

		if intf != nil {
			if len(field["1"]) == 0 {
				fmt.Printf("field %s does not contain any values. Skipping\n", tag)
				continue
			}
			if err := json.Unmarshal(*field["1"][0], intf); err != nil {
				return nil, fmt.Errorf("error de-marshalling value %s: %v", tag, err)
			}
			n, err := intf.Nodify(tag, m)
			if err != nil {
				return nil, fmt.Errorf("error nodifying value %s: %w", tag, err)
			}

			m.EnsureHierarchy(model.TagHierarchy, root)
			categoryTag := model.TagHierarchy[len(model.TagHierarchy)-1]
			category := root.FindCategory(categoryTag)
			if category == nil {
				return nil, fmt.Errorf("could not find category %d in tree", categoryTag)
			}
			category.AddChild(n)
		}

	}

	return root, nil
}
