package fsm

type ApplyResponse struct {
	Error error       `json:"err"`
	Data  interface{} `json:"data"`
}
