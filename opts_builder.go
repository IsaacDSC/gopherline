package workqueue

// OptsBuilder provides a fluent interface for building Opts
type OptsBuilder struct {
	opts Opts
}

// NewOptsBuilder creates a new OptsBuilder instance
func NewOptsBuilder() *OptsBuilder {
	return &OptsBuilder{opts: Opts{}}
}

// WithQueueType sets the queue type
func (b *OptsBuilder) WithQueueType(queueType string) *OptsBuilder {
	b.opts.QueueType = queueType
	return b
}

// WithMaxRetries sets the maximum number of retries
func (b *OptsBuilder) WithMaxRetries(maxRetries uint) *OptsBuilder {
	b.opts.MaxRetries = maxRetries
	return b
}

// WithScheduleIn sets the schedule in duration
func (b *OptsBuilder) WithScheduleIn(scheduleIn Duration) *OptsBuilder {
	b.opts.ScheduleIn = scheduleIn
	return b
}

// WithRetention sets the retention duration
func (b *OptsBuilder) WithRetention(retention Duration) *OptsBuilder {
	b.opts.Retention = retention
	return b
}

// WithUniqueTTL sets the unique TTL duration
func (b *OptsBuilder) WithUniqueTTL(uniqueTTL Duration) *OptsBuilder {
	b.opts.UniqueTTL = uniqueTTL
	return b
}

// Build returns the constructed Opts instance
func (b *OptsBuilder) Build() Opts {
	return b.opts
}
