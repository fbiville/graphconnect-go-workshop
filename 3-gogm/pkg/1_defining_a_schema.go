package workshop_gogm

import (
	"fmt"
	"github.com/mindstand/gogm/v2"
	"reflect"
)

// Package schema contains the gogm schema adapted from 3_result_mapping in section 2

// Person defines a person and their relationships for the schema
type Person struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	// Relationships
	Projects []*WorksOnEdge `gogm:"direction=outgoing;relationship=WORKS_ON"`
}

func NewPerson(name string) *Person {
	// note that we are initializing the relationship slice
	return &Person{
		Name:     name,
		Projects: []*WorksOnEdge{},
	}
}

// Topic defines a topic and their relationships for the schema
type Topic struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	// Relationships
	Projects []*Project `gogm:"direction=incoming;relationship=RELATES_TO"`
}

func NewTopic(name string) *Topic {
	// note that we are initializing the relationship slice
	return &Topic{
		Name:     name,
		Projects: []*Project{},
	}
}

// Project defines a project and its relationships
type Project struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	Type string `gogm:"name=project_type"`
	// Relationships
	Topics []*Topic       `gogm:"direction=outgoing;relationship=RELATES_TO"`
	People []*WorksOnEdge `gogm:"direction=incoming;relationship=WORKS_ON"`
}

func NewProject(name, _type string) *Project {
	// note that we are initializing the relationship slice
	return &Project{
		Name:   name,
		Type:   _type,
		Topics: []*Topic{},
		People: []*WorksOnEdge{},
	}
}

// WorksOnEdge implements gogm.Edge
// This will be phased out in a near future update for struct tags
type WorksOnEdge struct {
	gogm.BaseUUIDNode
	Start *Person
	End   *Project

	Role string `gogm:"name=name"`
}

func NewWorksOnEdge(start *Person, end *Project, role string) *WorksOnEdge {
	return &WorksOnEdge{
		Start: start,
		End:   end,
		Role:  role,
	}
}

func (w *WorksOnEdge) GetStartNode() interface{} {
	return w.Start
}

func (w *WorksOnEdge) GetStartNodeType() reflect.Type {
	return reflect.TypeOf(&Person{})
}

func (w *WorksOnEdge) SetStartNode(v interface{}) error {
	_start, ok := v.(*Person)
	if !ok {
		return fmt.Errorf("could not cast %T to *Person", v)
	}
	w.Start = _start
	return nil
}

func (w *WorksOnEdge) GetEndNode() interface{} {
	return w.End
}

func (w *WorksOnEdge) GetEndNodeType() reflect.Type {
	return reflect.TypeOf(&Project{})
}

func (w *WorksOnEdge) SetEndNode(v interface{}) error {
	_end, ok := v.(*Project)
	if !ok {
		return fmt.Errorf("could not cast %T to *Project", v)
	}
	w.End = _end
	return nil
}
