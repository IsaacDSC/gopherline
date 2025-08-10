package gopherline

type Opts struct {
	QueueType  string   `json:"queue_type,omitempty"`
	MaxRetries uint     `json:"max_retries,omitempty"`
	ScheduleIn Duration `json:"schedule_in,omitempty"`
	Retention  Duration `json:"retention,omitempty"`
	UniqueTTL  Duration `json:"unique_ttl,omitempty"`
}
