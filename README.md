# gopherline

A Go library for publishing events to a queue system with flexible configuration options, retry mechanisms, and scheduling capabilities.

## Features

- 🚀 Simple and intuitive event publishing API
- ⚙️ Flexible configuration with builder patterns
- 🔄 Configurable retry mechanisms
- ��� Event scheduling and retention control
- 🏷️ Queue type classification (critical, medium, low priority)
- 📊 Built-in correlation and event ID tracking
- 🌐 HTTP-based event publishing
- 📝 JSON serialization support

## Installation

```bash
go get github.com/IsaacDSC/gopherline
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/IsaacDSC/gopherline"
    "github.com/IsaacDSC/gopherline/SDK"
)

func main() {
    ctx := context.Background()
    
    // Create a producer with default options
    producer := SDK.NewProducer("http://localhost:8080", "your-token", gopherline.Opts{})
    
    // Build and publish an event
    payload := gopherline.NewInputBuilder().
        WithEvent("user.created").
        WithData(map[string]any{"user_id": "123", "email": "user@example.com"}).
        Build()
    
    if err := producer.Publish(ctx, payload); err != nil {
        panic(err)
    }
}
```

## Configuration

### Producer Configuration

Create a producer with custom default options:

```go
// Configure default options for all events
opts := gopherline.NewOptsBuilder().
    WithQueueType("internal.critical").
    WithMaxRetries(5).
    WithRetention(gopherline.NewDuration("168h")).  // 7 days
    WithScheduleIn(gopherline.NewDuration("5m")).   // Schedule 5 minutes from now
    Build()

producer := SDK.NewProducer("http://localhost:8080", "your-token", opts)
```

### Event-Specific Options

Override default options for specific events:

```go
// Create event-specific options
eventOpts := gopherline.NewOptsBuilder().
    WithQueueType("internal.high").
    WithMaxRetries(3).
    WithUniqueTTL(gopherline.NewDuration("1h")).
    Build()

payload := gopherline.NewInputBuilder().
    WithEvent("payment.processed").
    WithData(paymentData).
    WithOptions(eventOpts).  // Override default options
    Build()
```

## API Reference

### Producer Interface

```go
type Producer interface {
    Publish(ctx context.Context, payload Input) error
}
```

### Input Builder

The `InputBuilder` provides a fluent interface for constructing event payloads:

```go
builder := gopherline.NewInputBuilder()

// Required fields
builder.WithEvent("event.name")        // Event name/type
builder.WithData(data)                 // Event payload data

// Optional fields
builder.WithOptions(opts)              // Event-specific options
builder.WithCorrelationID("corr-123")  // For request tracking
builder.WithEventID("evt-456")         // Unique event identifier

input := builder.Build()
```

### Options Builder

Configure event processing behavior:

```go
opts := gopherline.NewOptsBuilder().
    WithQueueType("internal.critical").     // Queue priority classification
    WithMaxRetries(5).                      // Maximum retry attempts
    WithScheduleIn(gopherline.NewDuration("10m")).  // Delay before processing
    WithRetention(gopherline.NewDuration("24h")).   // How long to keep the event
    WithUniqueTTL(gopherline.NewDuration("1h")).    // Deduplication window
    Build()
```

### Duration Helper

Use human-readable duration strings:

```go
// Valid duration formats
gopherline.NewDuration("30s")    // 30 seconds
gopherline.NewDuration("5m")     // 5 minutes
gopherline.NewDuration("2h")     // 2 hours
gopherline.NewDuration("7d")     // 7 days (parsed as 168h)
gopherline.NewDuration("1h30m")  // 1 hour 30 minutes
```

## Configuration Options

### Queue Types

Classify events by priority or processing requirements:

- `internal.critical` - High priority, immediate processing
- `internal.high` - High priority
- `internal.medium` - Standard priority
- `internal.low` - Low priority, batch processing

### Retry Configuration

- `MaxRetries`: Number of retry attempts (0 = no retries)
- Exponential backoff is typically handled by the server

### Scheduling Options

- `ScheduleIn`: Delay before first processing attempt
- `Retention`: How long to keep the event in the system
- `UniqueTTL`: Deduplication window for identical events

## Usage Examples

### Basic Event Publishing

```go
func publishUserEvent(producer gopherline.Producer, userID string) error {
    ctx := context.Background()
    
    payload := gopherline.NewInputBuilder().
        WithEvent("user.updated").
        WithData(map[string]any{
            "user_id": userID,
            "timestamp": time.Now(),
        }).
        Build()
    
    return producer.Publish(ctx, payload)
}
```

### Critical Event with Retries

```go
func publishPaymentEvent(producer gopherline.Producer, payment Payment) error {
    ctx := context.Background()
    
    opts := gopherline.NewOptsBuilder().
        WithQueueType("internal.critical").
        WithMaxRetries(10).
        WithRetention(gopherline.NewDuration("72h")).
        Build()
    
    payload := gopherline.NewInputBuilder().
        WithEvent("payment.processed").
        WithData(payment).
        WithOptions(opts).
        WithCorrelationID(payment.TransactionID).
        Build()
    
    return producer.Publish(ctx, payload)
}
```

### Scheduled Event

```go
func scheduleReminder(producer gopherline.Producer, reminder Reminder) error {
    ctx := context.Background()
    
    opts := gopherline.NewOptsBuilder().
        WithQueueType("internal.medium").
        WithScheduleIn(gopherline.NewDuration("24h")). // Send in 24 hours
        WithMaxRetries(3).
        Build()
    
    payload := gopherline.NewInputBuilder().
        WithEvent("reminder.send").
        WithData(reminder).
        WithOptions(opts).
        Build()
    
    return producer.Publish(ctx, payload)
}
```

### Batch Processing with Deduplication

```go
func publishBatchEvents(producer gopherline.Producer, events []Event) error {
    ctx := context.Background()
    
    opts := gopherline.NewOptsBuilder().
        WithQueueType("internal.low").
        WithUniqueTTL(gopherline.NewDuration("1h")). // Deduplicate within 1 hour
        WithMaxRetries(2).
        Build()
    
    for _, event := range events {
        payload := gopherline.NewInputBuilder().
            WithEvent("batch.process").
            WithData(event).
            WithOptions(opts).
            WithEventID(fmt.Sprintf("batch-%s", event.ID)). // For deduplication
            Build()
        
        if err := producer.Publish(ctx, payload); err != nil {
            return fmt.Errorf("failed to publish event %s: %w", event.ID, err)
        }
    }
    
    return nil
}
```

### Service Integration Pattern

```go
type EventService struct {
    producer gopherline.Producer
}

func NewEventService(endpoint, token string) *EventService {
    opts := gopherline.NewOptsBuilder().
        WithQueueType("internal.medium").
        WithMaxRetries(3).
        WithRetention(gopherline.NewDuration("48h")).
        Build()
    
    producer := SDK.NewProducer(endpoint, token, opts)
    
    return &EventService{producer: producer}
}

func (s *EventService) PublishUserCreated(ctx context.Context, user User) error {
    payload := gopherline.NewInputBuilder().
        WithEvent("user.created").
        WithData(user).
        WithCorrelationID(user.ID).
        Build()
    
    return s.producer.Publish(ctx, payload)
}

func (s *EventService) PublishCriticalAlert(ctx context.Context, alert Alert) error {
    opts := gopherline.NewOptsBuilder().
        WithQueueType("internal.critical").
        WithMaxRetries(10).
        Build()
    
    payload := gopherline.NewInputBuilder().
        WithEvent("alert.critical").
        WithData(alert).
        WithOptions(opts). // Override default options
        Build()
    
    return s.producer.Publish(ctx, payload)
}
```

## Error Handling

The library returns detailed errors for various failure scenarios:

```go
if err := producer.Publish(ctx, payload); err != nil {
    switch {
    case strings.Contains(err.Error(), "event cannot be empty"):
        // Handle validation error
        log.Printf("Invalid event: %v", err)
    case strings.Contains(err.Error(), "failed to marshal"):
        // Handle serialization error
        log.Printf("Serialization failed: %v", err)
    case strings.Contains(err.Error(), "failed on publisher"):
        // Handle network/HTTP error
        log.Printf("Publishing failed: %v", err)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

## Best Practices

### 1. Use Appropriate Queue Types
- `critical`: Payment processing, security alerts
- `high`: User notifications, important updates  
- `medium`: General business events, data sync
- `low`: Analytics, logging, batch operations

### 2. Set Reasonable Retry Limits
- Critical events: 5-10 retries
- Standard events: 3-5 retries
- Analytics/logging: 1-2 retries

### 3. Configure Retention Appropriately
- Compliance events: Long retention (weeks/months)
- Notifications: Short retention (hours/days)
- Analytics: Medium retention (days/weeks)

### 4. Use Correlation IDs
Always include correlation IDs for request tracing:

```go
payload := gopherline.NewInputBuilder().
    WithEvent("order.created").
    WithData(order).
    WithCorrelationID(requestID). // From incoming request
    Build()
```

### 5. Handle Context Cancellation
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := producer.Publish(ctx, payload); err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Publishing timed out")
    }
    return err
}
```

## Server Configuration

The producer connects to a gopherline server endpoint. Ensure your server:

1. Accepts POST requests at `/event/publisher`
2. Uses Basic authentication with the provided token
3. Accepts JSON payloads with the expected structure
4. Returns appropriate HTTP status codes

## License

This project is licensed under the terms specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
