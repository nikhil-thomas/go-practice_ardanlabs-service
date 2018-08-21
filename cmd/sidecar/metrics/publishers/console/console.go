package console

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Collector defines a contract a collector must support
// so a consume can retrieve metrics
type Collector interface {
	Collect() (map[string]interface{}, error)
}

// Console provides the ability to receive metrics
// from internal services using expvar
type Console struct {
	collector Collector
	wg        sync.WaitGroup
	timer     *time.Timer
	shutdown  chan struct{}
}

// New creates a console based consumer
func New(collector Collector, interval time.Duration) (*Console, error) {
	con := Console{
		collector: collector,
		timer:     time.NewTimer(interval),
		shutdown:  make(chan struct{}),
	}

	con.wg.Add(1)
	go func() {
		defer con.wg.Done()
		for {
			con.timer.Reset(interval)
			select {
			case <-con.timer.C:
				con.publish()
			case <-con.shutdown:
				return
			}
		}
	}()

	return &con, nil
}

// Stop stops the goroutine whiich collects metrics
func (con *Console) Stop() {
	close(con.shutdown)
	con.wg.Wait()
}

// publish pulls the metrics and publishes them to the console
func (con *Console) publish() {
	data, err := con.collector.Collect()
	if err != nil {
		log.Println(err)
		return
	}

	out, err := marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(out)
}

// marshal marshals data from map to json string
func marshal(data map[string]interface{}) (string, error) {
	out, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(out), nil
}
