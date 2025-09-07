package workqueue

// InputBuilder provides a fluent interface for building Input
type InputBuilder struct {
	input Input
}

// NewInputBuilder creates a new InputBuilder instance
func NewInputBuilder() *InputBuilder {
	return &InputBuilder{input: Input{}}
}

// WithServiceName sets the service name
func (b *InputBuilder) WithServiceName(serviceName string) *InputBuilder {
	b.input.ServiceName = serviceName
	return b
}

// WithEvent sets the event name
func (b *InputBuilder) WithEvent(event string) *InputBuilder {
	b.input.Event = event
	return b
}

// WithData sets the data
func (b *InputBuilder) WithData(data any) *InputBuilder {
	b.input.Data = data
	return b
}

// WithOptions sets the options map
func (b *InputBuilder) WithOptions(options Opts) *InputBuilder {
	b.input.Options = options
	return b
}

// WithCorrelationID sets the correlation ID
func (b *InputBuilder) WithCorrelationID(correlationID string) *InputBuilder {
	b.input.CorrelationID = correlationID
	return b
}

// WithEventID sets the event ID
func (b *InputBuilder) WithEventID(eventID string) *InputBuilder {
	b.input.EventID = eventID
	return b
}

// Build returns the constructed Input instance
func (b *InputBuilder) Build() Input {
	return b.input
}
