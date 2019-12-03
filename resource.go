package main

import "fmt"

// Resource is a resource
type Resource struct {
	ID    string
	Name  string
	Desc  string
	Value int
	Max   int
	Min   int
	Delta int
}

// NewResource creates a resource
func NewResource(desc, name string, value, delta int) *Resource {
	return &Resource{NewID("resource"), name, desc, value, value, 0, delta}
}

// Update will apply the delta within the  min/max bounds
func (r *Resource) Update() {
	v := r.Value + r.Delta
	if v > r.Max {
		v = r.Max
	} else if v < r.Min {
		v = r.Min
	}
	r.Value = v
}

// String returns a snapshot in format name:value/max
func (r *Resource) String() string {
	return fmt.Sprintf("%s:%v/%v", r.Name, r.Value, r.Max)
}
