// Package schema contains the gogm schema identical to the schema used in 3_result_mapping in section 2
package schema

import (
	"github.com/mindstand/gogm/v2"
)

// Person defines a person and their relationships for the schema
type Person struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	// Relationships
	Projects []*Project `gogm:"direction=outgoing;relationship=WORKS_ON"`
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

// Project defines a project and its relationships
type Project struct {
	// Required (or use gogm.BaseNode)
	gogm.BaseUUIDNode
	// Members
	Name string `gogm:"name=name;index"`
	Type string `gogm:"name=project_type"`
	// Relationships
	Topics []*Topic  `gogm:"direction=outgoing;relationship=RELATES_TO"`
	People []*Person `gogm:"direction=incoming;relationship=WORKS_ON"`
}
