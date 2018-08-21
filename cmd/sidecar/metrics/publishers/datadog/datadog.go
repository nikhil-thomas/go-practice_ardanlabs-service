package datadog

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Collector defines a collector interface
// to retrieve metrics
type Collector interface {
	Collect() (map[string]interface{}, error)
}

// Datadog provides te ability to receive metrics
// from internal services using expvar
type Datadog struct {
	collector Collector
	wg        sync.WaitGroup
	timer     *time.Timer
	shutdown  chan struct{}
}

// New creates a Datadog based consumer
func New(collector Collector, interval time.Duration) (*Datadog, error) {
	dg := Datadog{
		collector: collector,
		timer:     time.NewTimer(interval),
		shutdown:  make(chan struct{}),
	}

	dg.wg.Add(1)
	go func() {
		defer dg.wg.Done()
		for {
			dg.timer.Reset(interval)
			select {
			case <-dg.timer.C:
				dg.publish()
			case <-dg.shutdown:
				return
			}
		}
	}()

	return &dg, nil
}

// Stop is used to shutdown the goroutine colelcting metrics
func (dg *Datadog) Stop() {
	close(dg.shutdown)
	dg.wg.Wait()
}

// publish collects metrics and publishes metrics in datadog
func (dg *Datadog) publish() {
	data, err := dg.collector.Collect()
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

// marshal handles the marshaling of the map to a datadog json
func marshal(data map[string]interface{}) (string, error) {

	/*
	   {"series" : [
	                   {
	                       "metric": "test.metric",
	                       "points": [
	                           [
	                               $currenttime,
	                               20
	                           ]
	                       ],
	                       "type": "guage",
	                       "host": "test.example.com",
	                       "tags": [
	                           "environment.test"
	                       ]
	                   }
	       ]
	   }
	*/

	// Extract base key values
	mType := "guage"
	host, ok := data["host"].(string)
	if !ok {
		host = "unknown"
	}
	env := "dev"
	if host != "localhost" {
		env = "prod"
	}
	envTag := "environment:" + env

	// define Datdog format
	type series struct {
		Metric string          `json:"metric"`
		Points [][]interface{} `json:"points"`
		Type   string          `json:"type"`
		Host   string          `json:"host"`
		Tags   []string        `json:"tags"`
	}

	type dog struct {
		Series []series `json:"series"`
	}

	// Populate the data into the data structure
	var d dog
	for key, value := range data {
		switch value.(type) {
		case int, float64:
			d.Series = append(d.Series, series{
				Metric: env + "." + key,
				Points: [][]interface{}{[]interface{}{"$currenttime", value}},
				Type:   mType,
				Host:   host,
				Tags:   []string{envTag},
			})
		}
	}

	out, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}
