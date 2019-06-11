package usecase

import (
	"encoding/json"
	"strings"
)

type (
	MatrixRenderer struct {
		Pipeline string
		Matrix   map[string][]string
	}

	IncludeRenderer struct {
		Pipeline string
		Include  []map[string]string
	}
)

func (pr *MatrixRenderer) ArgName() string {
	if pr.Matrix == nil {
		return ""
	}
	a := make([]string, len(pr.Matrix))
	i := 0
	for k := range pr.Matrix {
		a[i] = k
		i++
	}
	return strings.Join(a, ", ")
}

func (pr *IncludeRenderer) ArgName() string {
	if len(pr.Include) == 0 {
		return ""
	}
	a := make([]string, len(pr.Include[0]))
	i := 0
	for k := range pr.Include[0] {
		a[i] = k
		i++
	}
	return strings.Join(a, ", ")
}

func (pr *IncludeRenderer) Args() string {
	b, _ := json.MarshalIndent(pr.Include, "", "  ")
	return string(b)
}

func (pr *IncludeRenderer) ArgNameWithArg() string {
	if len(pr.Include) == 0 {
		return ""
	}
	a := make([]string, len(pr.Include[0]))
	i := 0
	for k := range pr.Include[0] {
		a[i] = "arg." + k
		i++
	}
	return strings.Join(a, ", ")
}

func (pr *MatrixRenderer) ArrayArgs() map[string]string {
	a := make(map[string]string, len(pr.Matrix))
	for k, v := range pr.Matrix {
		b, _ := json.MarshalIndent(v, "", "  ")
		a[k] = string(b)
	}
	return a
}
