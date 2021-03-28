package shell_test

import (
	"testing"

	"github.com/kajjagtenberg/eventflowdb/shell"
	"github.com/stretchr/testify/assert"
)

func TestBabel(t *testing.T) {
	assert := assert.New(t)

	babel := shell.NewBabel()

	compiled, err := babel.Compile("const pow = x => x*x;")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(compiled), 61) // Not the best test
}
