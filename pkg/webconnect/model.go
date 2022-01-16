package webconnect

import (
	"fmt"
)

type Model map[string]*ObjectMeta

type ObjectMeta struct {
	Prio              int      `json:"Prio"`
	TagId             int      `json:"TagId"`
	TagIdEventMessage int      `json:"TagIdEvtMsg"`
	Unit              *int     `json:"Unit"`
	DataFormat        int      `json:"DataFrmt"`
	Scale             *float64 `json:"Scale"`
	Type              int      `json:"Type"`
	WriteLevel        int      `json:"WriteLevel"`
	TagHierarchy      []int    `json:"TagHier"`
	Min               *bool    `json:"Min"`
	Max               *bool    `json:"Max"`
	Average           *bool    `json:"Avg"`
	Count             *bool    `json:"Cnt"`
	MinD              *bool    `json:"MinD"`
	MaxD              *bool    `json:"MaxD"`
}

type Language map[int]string

type Meta struct {
	model    *Model
	language *Language
}

func (m *Meta) GetTranslation(tag int) (*string, error) {
	text, ok := (*m.language)[tag]
	if !ok {
		return nil, fmt.Errorf("tag %s not found in language definition", tag)
	}
	return &text, nil
}

func (m *Meta) GetModel(tag string) (*ObjectMeta, error) {
	model, ok := (*m.model)[tag]
	if !ok {
		return nil, fmt.Errorf("tag %s not found in model definition", tag)
	}
	return model, nil
}
