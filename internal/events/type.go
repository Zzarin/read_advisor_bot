package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct { //type Event is common struct for all messengers, so we can't add chatID and userName fields from telegram.
	Type Type
	Text string
	Meta interface{} //we define an empty interface(not good), so we can add Meta struct with chatID and userName in the Event struct.
}
