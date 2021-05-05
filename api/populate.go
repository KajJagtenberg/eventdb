package api

import (
	"encoding/json"
	"sync"

	"github.com/KajJagtenberg/eventflowdb/store"
	"github.com/google/uuid"
)

type PopulateRequest struct {
	Count int `json:"count"`
	Size  int `json:"size"`
}

func Populate(s store.Store, c *Ctx) error {
	if len(c.Args) == 0 {
		return ErrInsufficientArguments
	}

	var req PopulateRequest
	if err := json.Unmarshal(c.Args[0], &req); err != nil {
		return err
	}

	if req.Count <= 0 {
		req.Count = 100
	}

	if req.Size <= 0 {
		req.Size = 100
	}

	wg := sync.WaitGroup{}

	for i := 0; i < req.Count; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			stream := uuid.New()

			// data := make([]byte, req.Size)
			// rand.Read(data)

			events := []store.EventData{
				{
					Type: "RandomEvent",
					Data: []byte(`{"result":[{"id":"625eecb62ae26ff883aa87e25a332941","status":"accepted","permissions":{"access":{"read":true,"edit":true},"analytics":{"read":true,"edit":false},"app":{"read":false,"edit":true},"auditlogs":{"read":true,"edit":false},"billing":{"read":true,"edit":true},"cache_purge":{"read":false,"edit":true},"dns_records":{"read":true,"edit":true},"lb":{"read":true,"edit":true},"legal":{"read":true,"edit":true},"logs":{"read":true,"edit":true},"member":{"read":true,"edit":true},"organization":{"read":true,"edit":true},"ssl":{"read":true,"edit":true},"stream":{"read":true,"edit":true},"subscription":{"read":true,"edit":true},"teams":{"read":true,"edit":true},"waf":{"read":true,"edit":true},"webhooks":{"read":true,"edit":true},"worker":{"read":true,"edit":true},"zone":{"read":true,"edit":true},"zone_settings":{"read":true,"edit":true}},"roles":["Super Administrator - All Privileges"],"account":{"id":"c784bdaf1c9280a281799c54729ad71a","name":"Info@kajjagtenberg.nl's Account","settings":{"enforce_twofactor":false,"access_approval_expiry":null},"meta":{"has_pro_zones":false,"has_business_zones":false,"has_enterprise_zones":false},"legacy_flags":{"railgun":{"enabled":false},"dns_firewall":{"enabled":false},"china_network_visible":{"enabled":false},"china_private_key_network_deployment":{"enabled":false},"custom_pages":{"enabled":false},"enterprise_zone_quota":{"maximum":0,"current":0,"available":0}}}}],"result_info":{"page":1,"per_page":20,"total_pages":1,"count":1,"total_count":1},"success":true,"errors":[],"messages":[]}`),
				},
			}

			s.Add(stream, 0, events)
		}()
	}

	wg.Wait()

	c.Conn.WriteString("OK")

	return nil
}
