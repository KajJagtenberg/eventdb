package tests

import (
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"go.etcd.io/bbolt"
)

func TestAdd(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	db, err := bbolt.Open(file.Name(), 0666, bbolt.DefaultOptions)
	p, err := persistence.NewPersistence(db)
	if err != nil {
		t.Fatal(err)
	}

	count := 1000

	total := 0

	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func() {
			stream := uuid.New()

			var events []persistence.EventData

			for j := 0; j < rand.Intn(5)+1; j++ {
				data := make([]byte, 100)
				metadata := make([]byte, 50)

				rand.Read(data)
				rand.Read(metadata)

				events = append(events, persistence.EventData{
					Type:     "Random",
					Data:     data,
					Metadata: metadata,
					AddedAt:  time.Now(),
				})

				total++
			}

			if _, err := p.Add(stream, 0, events); err != nil {
				t.Fatal(err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	t.Log(total)
}
