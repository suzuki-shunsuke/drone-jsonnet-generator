package domain

const (
	MatrixTemplate = `local pipeline({{ .ArgName }}) = {{ .Pipeline }};

{{range $k, $v := .ArrayArgs -}}
local array_{{ $k }} = {{ $v }};
{{end}}
[
  pipeline({{ .ArgName }}) {{range $k, $v := .ArrayArgs }}for {{ $k }} in array_{{ $k }} {{end}}
]
`

	IncludeTemplate = `local pipeline({{ .ArgName }}) = {{ .Pipeline }};

local args = {{ .Args }};

[
  pipeline({{ .ArgNameWithArg }}) for arg in args
]
`
)
