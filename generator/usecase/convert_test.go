package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/domain"
)

func Test_validateArg(t *testing.T) {
	data := []struct {
		title string
		arg   *domain.ConvertArg
		isErr bool
	}{
		{
			title: "arg is nil",
			isErr: true,
		},
		{
			title: "source is required",
			arg: &domain.ConvertArg{
				Stdout: true,
			},
			isErr: true,
		},
		{
			title: "target or stdout is required",
			arg: &domain.ConvertArg{
				Source: ".drone.yml",
			},
			isErr: true,
		},
		{
			title: "source and target are set",
			arg: &domain.ConvertArg{
				Source: ".drone.yml",
				Target: ".drone.jsonnet",
			},
			isErr: false,
		},
		{
			title: "source and stdout are set",
			arg: &domain.ConvertArg{
				Source: ".drone.yml",
				Stdout: true,
			},
			isErr: false,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := validateArg(d.arg)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
		})
	}
}
