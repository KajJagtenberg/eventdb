package fsm

type CommandPayload struct {
	Operation string      `json:"op"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}
