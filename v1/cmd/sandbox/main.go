package main

import (
	store "eventflowdb/store/proto"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func main() {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
	id, err := ulid.MustNew(ulid.Now(), entropy).MarshalBinary()

	stream, err := uuid.New().MarshalBinary()
	if err != nil {
		log.Fatalf("Unable to encode stream: %v", err)
	}

	m := &store.RecordedEvent{
		Id:            id,
		Stream:        stream,
		Version:       0,
		Type:          "ProductAdded",
		Data:          []byte(`{"name":"Samsung Galaxy S8","price":80000}`),
		CausationId:   id,
		CorrelationId: id,
		AddedAt:       time.Now().Unix(),
	}

	packed, err := proto.Marshal(m)
	if err != nil {
		log.Fatalf("Unable to marshal: %v", err)
	}

	log.Println(len(packed))
}
