package consumer

import (
	"log"
	"read_advisor_bot/internal/events"
	"sync"
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

func (c *Consumer) handleEvents(events []events.Event) error {
	var wg sync.WaitGroup
	for _, event := range events {
		wg.Add(1)
		event := event
		go func() {
			defer wg.Done()
			log.Printf("got new event: %s", event.Text)

			err := c.processor.Process(event)
			if err != nil {
				log.Printf("got new event: %s", err.Error())
			}
		}()
	}
	wg.Wait()
	return nil
}
