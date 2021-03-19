package store_test

import (
	protos "eventflowdb/gen"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestStreamDeserialization(t *testing.T) {
	var data []byte

	var m protos.Stream

	if err := proto.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	t.Log(m)
}
