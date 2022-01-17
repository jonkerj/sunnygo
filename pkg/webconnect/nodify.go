package webconnect

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jonkerj/sunnygo/pkg/tree"
)

type Nodifyable interface {
	nodify(string, *Meta) (tree.Node, error)
}

func (m *Meta) ensureHierarchy(hierarchy []int, root *tree.CategoryNode) (*tree.CategoryNode, error) {
	parent := root
	for _, categoryTag := range hierarchy {
		r := parent.GetCategory(categoryTag)
		if r != nil {
			parent = r
			continue
		} else {
			// this category is not yet present in the tree. Create the node
			translation, err := m.GetTranslation(categoryTag)
			if err != nil {
				return nil, fmt.Errorf("could not complete category tree due to missing translation: %v", err)
			}

			newCategory := tree.NewCategoryNode(categoryTag, *translation)
			parent.AddChild(newCategory)
			parent = newCategory
		}

	}

	return parent, nil
}

func (i *IntValue) nodify(tag string, meta *Meta) (tree.Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	n := tree.NewValueNode(tag, *name)

	if model.Unit != nil {
		unit, err := meta.GetTranslation(*model.Unit)
		if err != nil {
			return nil, fmt.Errorf("could not translation for unit %d: %w", *model.Unit, err)
		}
		n.SetUnit(unit)
	}

	n.SetValue(float64(i.Val) * (*model.Scale))

	return n, nil
}

func (s *StringValue) nodify(tag string, meta *Meta) (tree.Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	n := tree.NewTextNode(tag, *name, s.Val)
	return n, nil
}

func (d *DurationValue) nodify(tag string, meta *Meta) (tree.Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	n := tree.NewDurationNode(tag, *name, time.Duration(d.Val)*time.Second)

	return n, nil
}

func (t *TagListValue) nodify(tag string, meta *Meta) (tree.Node, error) {
	model, err := meta.GetModel(tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	name, err := meta.GetTranslation(model.TagId)
	if err != nil {
		return nil, fmt.Errorf("could not translation for tag %d: %w", model.TagId, err)
	}

	if len(t.Val) == 0 {
		n := tree.NewTextNode(tag, *name, "n/a")
		return n, nil
	}

	text, err := meta.GetTranslation(t.Val[0].Tag)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve model for tag %s: %w", tag, err)
	}

	n := tree.NewTextNode(tag, *name, *text)

	return n, nil
}

func NodifyAllValues(deviceId string, m *Meta, response *ResultReponse) (*tree.CategoryNode, error) {
	fields, ok := response.Result[deviceId]
	if !ok {
		return nil, fmt.Errorf("device %s not found in values", deviceId)
	}

	rootName := "root"
	root := tree.NewCategoryNode(0, rootName)

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
			n, err := intf.nodify(tag, m)
			if err != nil {
				return nil, fmt.Errorf("error nodifying value %s: %w", tag, err)
			}

			category, err := m.ensureHierarchy(model.TagHierarchy, root)
			if err != nil {
				return nil, fmt.Errorf("error retrieving category node: %w", err)
			}
			if category == nil {
				return nil, fmt.Errorf("could not create hierarchy %s in tree", model.TagHierarchy)
			}
			category.AddChild(n)
		}

	}

	return root, nil
}
