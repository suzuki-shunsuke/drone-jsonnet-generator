package domain

type (
	ConvertArg struct {
		Source string
		Target string
		Stdout bool
	}

	LegacyYAML struct {
		Matrix map[string][]string
	}

	Include struct {
		Include []map[string]string
	}

	LegacyIncludedYAML struct {
		Matrix Include
	}
)
