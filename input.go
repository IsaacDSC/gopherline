package workqueue

type Input struct {
	Event   string `json:"event_name"`
	Data    any    `json:"data"`
	Options Opts   `json:"opts"`
	//These fields are used to pass metadata.headers
	CorrelationID string `json:"correlation_id"`
	EventID       string `json:"event_id"`
}
