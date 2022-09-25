package events

import "context"

// Fetcher интерфейс под который необходимо реализовать event-структуру клиента.
// Здесь он реализован для telegram.
type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

// Processor интерфейс под который необходимо реализовать event-структуру клиента.
// Здесь он реализован для telegram.
type Processor interface {
	Process(ctx context.Context, e Event) error
}

type Storage interface {
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
