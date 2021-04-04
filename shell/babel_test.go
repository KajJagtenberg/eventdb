package shell_test

import (
	"strings"
	"testing"

	"github.com/kajjagtenberg/eventflowdb/shell"
	"github.com/stretchr/testify/assert"
)

func TestBabel(t *testing.T) {
	assert := assert.New(t)

	babel, err := shell.NewBabel()
	if err != nil {
		t.Fatal(err)
	}

	code, err := babel.Compile("const pow = x => x * x;")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(strings.Contains(code, `function pow(x)`), true)
}
