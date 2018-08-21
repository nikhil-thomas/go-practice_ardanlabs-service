package expvar

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Expvar provides the ability to receive metrics from internal services using expvar
type Expvar struct {
	host   string
	tr     *http.Transport
	client http.Client
}

// New creates a Expvar for collection metrics
func New(host string) (*Expvar, error) {
	tr := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   20 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          2,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	exp := Expvar{
		host: host,
		tr:   &tr,
		client: http.Client{
			Transport: &tr,
			Timeout:   1 * time.Second,
		},
	}

	return &exp, nil
}

// Collect pulls metrics from the configured host
// Collect implements the console.Colelctor interface
func (exp *Expvar) Collect() (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, exp.host, nil)
	log.Println(exp.host)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := exp.client.Do(req)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	memStats, ok := (data["memstats"]).(map[string]interface{})
	if ok {
		data["heap"] = memStats["Alloc"]
	}

	u, err := url.Parse(exp.host)
	if err != nil {
		return nil, err
	}
	data["host"] = u.Hostname()

	delete(data, "memStats")
	delete(data, "cmdline")

	return data, nil
}
