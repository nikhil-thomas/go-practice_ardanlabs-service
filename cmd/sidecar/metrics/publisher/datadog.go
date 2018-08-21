package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// Datadog provides the ability to publish metrics to Datadog
type Datadog struct {
	apiKey string
	host   string
	tr     *http.Transport
	client http.Client
}

// NewDatadog initializes Datadog access for publishing metrics
func NewDatadog(apiKey string, host string) *Datadog {
	tr := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          2,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	d := Datadog{
		apiKey: apiKey,
		host:   host,
		tr:     &tr,
		client: http.Client{
			Transport: &tr,
			Timeout:   1 * time.Second,
		},
	}

	return &d
}

// Publish handles the processing of metrics for delivery to Datadog
func (d *Datadog) Publish(data map[string]interface{}) {
	doc, err := marshalDatadog(data)
	if err != nil {
		log.Println("datadog.publish", err)
		return
	}

	if err := sendDatadog(d, doc); err != nil {
		log.Println("datadog.publish : published :", string(doc))
		return
	}
	log.Println("datadog.publish : published :", string(doc))
}

// marshalDatadog converts the data map to json
func marshalDatadog(data map[string]interface{}) ([]byte, error) {
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

	type doc struct {
		Series []series `json:"series"`
	}

	// Populate the data into the data structure
	var d doc
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
		return nil, err
	}

	return out, nil
}

// sendDatadog sends datato Datadog SAAS
func sendDatadog(d *Datadog, data []byte) error {
	url := fmt.Sprintf("%s?api_key=%s", d.host, d.apiKey)
	b := bytes.NewBuffer(data)

	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Printf("datadog.publish : error : status[%d] : %s", resp.StatusCode, out)
	}
	return nil
}
