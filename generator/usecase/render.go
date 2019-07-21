package usecase

import (
	"encoding/json"
	"strings"
)

type (
	MatrixRenderer struct {
		Pipeline string
		Matrix   []Axis
	}

	IncludeRenderer struct {
		Pipeline string
		Include  [][]IncludeAxis
	}

	Axis struct {
		Name   string
		Values []string
	}

	IncludeAxis struct {
		Name  string
		Value string
	}
)

func (pr *MatrixRenderer) ArgName() string {
	if pr.Matrix == nil {
		return ""
	}
	a := make([]string, len(pr.Matrix))
	for i, k := range pr.Matrix {
		a[i] = k.Name
	}
	return strings.Join(a, ", ")
}

func (pr *MatrixRenderer) ArrayArgs() map[string]string {
	a := make(map[string]string, len(pr.Matrix))
	for _, v := range pr.Matrix {
		b, _ := json.MarshalIndent(v.Values, "", "  ")
		a[v.Name] = string(b)
	}
	return a
}

func (pr *IncludeRenderer) ArgName() string {
	if len(pr.Include) == 0 {
		return ""
	}
	a := make([]string, len(pr.Include[0]))
	for i, k := range pr.Include[0] {
		a[i] = k.Name
	}
	return strings.Join(a, ", ")
}

func (pr *IncludeRenderer) Args() string {
	a := make([]map[string]string, len(pr.Include))
	for i, v := range pr.Include {
		m := make(map[string]string, len(v))
		for _, c := range v {
			m[c.Name] = c.Value
		}
		a[i] = m
	}
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func (pr *IncludeRenderer) ArgNameWithArg() string {
	if len(pr.Include) == 0 {
		return ""
	}
	a := make([]string, len(pr.Include[0]))
	i := 0
	for _, k := range pr.Include[0] {
		a[i] = "arg." + k.Name
		i++
	}
	return strings.Join(a, ", ")
}
