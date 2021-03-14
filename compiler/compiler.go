package compiler

import (
	_ "embed"
	"io/ioutil"
	"strings"

	babel "github.com/jvatic/goja-babel"
)

func Compile(source string) (string, error) {
	output, err := babel.Transform(strings.NewReader(source), map[string]interface{}{
		"presets": []string{"env"},
	})
	if err != nil {
		return "", err
	}

	code, err := ioutil.ReadAll(output)
	if err != nil {
		return "", err
	}

	return string(code), nil
}
