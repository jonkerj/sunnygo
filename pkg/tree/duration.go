package tree

import (
	"fmt"
	"time"
)

type DurationNode struct {
	tag      string
	name     string
	duration time.Duration
}

func NewDurationNode(tag string, name string, duration time.Duration) *DurationNode {
	d := &DurationNode{
		tag:      tag,
		name:     name,
		duration: duration,
	}

	return d
}

func (d DurationNode) Children() []Node {
	return []Node{}
}

func (d *DurationNode) String() string {
	return fmt.Sprintf("%s: %s", d.name, d.duration)
}
