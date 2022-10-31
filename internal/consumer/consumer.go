package consumer

import (
	"log"
	"read_advisor_bot/internal/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // how many events will be process by once
}

type Starter interface {
	Start() error
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		err = c.handleEvents(gotEvents)
		if err != nil {
			log.Printf("[ERR] consumer, handle events: %s", err.Error())
			continue
		}
	}
}

/*potential problems:
1. Event lost: try retry
2. Handling of all events: stop
after first error, error counter
3. Concurrent processing - homework. Need wait
group (sync package)
*/
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		err := c.processor.Process(event)
		if err != nil {
			log.Printf("got new event: %s", err.Error())
			continue
		}
	}
	return nil
}
