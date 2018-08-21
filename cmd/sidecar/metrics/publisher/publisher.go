package publisher

import (
	"log"
	"sync"
	"time"
)

// Collector defines a contract a collector must support
// so a consume can retrieve metrics
type Collector interface {
	Collect() (map[string]interface{}, error)
}

// Publisher defines a handler function that is invoked on each interval
type Publisher func(map[string]interface{})

// Publish provides the ability to receive metrics on an interval
type Publish struct {
	collector Collector
	publisher Publisher
	wg        sync.WaitGroup
	timer     *time.Timer
	shutdown  chan struct{}
}

// New creates a Publish for consuming and publishing metrics
func New(collector Collector, publisher Publisher, interval time.Duration) (*Publish, error) {
	p := Publish{
		collector: collector,
		publisher: publisher,
		timer:     time.NewTimer(interval),
		shutdown:  make(chan struct{}),
	}

	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		for {
			p.timer.Reset(interval)
			select {
			case <-p.timer.C:
				p.update()
			case <-p.shutdown:
				return
			}
		}
	}()

	return &p, nil
}

// Stop is used to shutdwon the goroutine colelcting metrics
func (p *Publish) Stop() {
	close(p.shutdown)
	p.wg.Wait()
}

// update pulls metrics and publishes to the specified system
func (p *Publish) update() {
	data, err := p.collector.Collect()
	if err != nil {
		log.Println(err)
		return
	}
	p.publisher(data)
}
