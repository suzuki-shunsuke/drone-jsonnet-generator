package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/drone/drone-yaml/yaml/converter/legacy"
	gyaml "github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/domain"
)

func validateArg(arg *domain.ConvertArg) error {
	if arg == nil {
		return errors.New("arg is nil")
	}
	if arg.Source == "" {
		return errors.New("source is required")
	}
	if !arg.Stdout && arg.Target == "" {
		return errors.New("target or stdout is required")
	}
	return nil
}

func Convert(arg *domain.ConvertArg) error {
	if err := validateArg(arg); err != nil {
		return err
	}
	f, err := os.Open(arg.Source)
	if err != nil {
		return errors.Wrap(err, "failed to open the source file .drone.yml "+arg.Source)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to read the source file .drone.yml "+arg.Source)
	}
	m, matrix, include, err := convertMatrix(b)
	if err != nil {
		return err
	}
	if matrix == nil && include == nil {
		// matrix isn't used
		return errors.New("matrix build isn't used")
	}

	// remove matrix and convert to byte
	delete(m, "matrix")
	bWithoutMatrix, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	// convert to v1 format
	yamlV1, err := legacy.Convert(bWithoutMatrix)
	if err != nil {
		return err
	}

	// convert yaml to json struct
	d := map[string]interface{}{}
	if err := gyaml.Unmarshal(yamlV1, &d); err != nil {
		return err
	}

	// convert pipeline name
	if matrix != nil {
		c := make([]string, len(matrix))
		i := 0
		for k := range matrix {
			c[i] = fmt.Sprintf("' %s:' + %s", k, k)
			i++
		}
		d["name"] = "'" + strings.Join(c, " + ")[2:]
	} else {
		c := make([]string, len(include[0]))
		i := 0
		for k := range include[0] {
			c[i] = fmt.Sprintf("' %s:' + %s", k, k)
			i++
		}
		d["name"] = "'" + strings.Join(c, " + ")[2:]
	}

	// convert YAML to JSON
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(d); err != nil {
		return errors.Wrap(err, "failed to convert YAML to JSON")
	}
	jsonV1 := strings.TrimSpace(buf.String())

	if arg.Stdout {
		return genJSONNet(os.Stdout, jsonV1, matrix, include)
	}
	w, err := os.Create(arg.Target)
	if err != nil {
		return err
	}
	defer w.Close()
	return genJSONNet(w, jsonV1, matrix, include)
}

func genJSONNet(w io.Writer, jsonV1 string, matrix map[string][]string, include []map[string]string) error {
	if include == nil {
		renderer := &MatrixRenderer{
			Pipeline: jsonV1,
			Matrix:   matrix,
		}
		tpl, err := template.New("matrix").Parse(domain.MatrixTemplate)
		if err != nil {
			return err
		}
		return tpl.Execute(w, renderer)
	}
	renderer := &IncludeRenderer{
		Pipeline: jsonV1,
		Include:  include,
	}
	tpl, err := template.New("include").Parse(domain.IncludeTemplate)
	if err != nil {
		return err
	}
	return tpl.Execute(w, renderer)
}

func convertMatrix(data []byte) (map[string]interface{}, map[string][]string, []map[string]string, error) {
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to parse the source .drone.yml")
	}
	if _, ok := m["matrix"]; !ok {
		return m, nil, nil, nil
	}

	legacyYAML := &domain.LegacyYAML{}
	err := yaml.Unmarshal(data, legacyYAML)
	if err == nil {
		return m, legacyYAML.Matrix, nil, nil
	}

	legacyIncludedYAML := &domain.LegacyIncludedYAML{}
	if err := yaml.Unmarshal(data, legacyIncludedYAML); err != nil {
		return m, nil, nil, errors.New("matrix format is invalid")
	}
	return m, nil, legacyIncludedYAML.Matrix.Include, nil
}
