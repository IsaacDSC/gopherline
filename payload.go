package workqueue

type Payload struct {
	ServiceName string         `json:"service_name"`
	Event       string         `json:"event_name"`
	Data        any            `json:"data"`
	Options     Opts           `json:"opts"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}
