package publisher

import (
	"encoding/json"
	"log"
)

// set possible publisher types
const (
	TypeConsole = "CONSOLE"
	TypeDatadog = "DATADOG"
)

// Console handles metrics for direct display on console
func Console(data map[string]interface{}) {
	out, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return
	}
	log.Println(string(out))
}

// Datadog processes metrics for delivry to Datadog
func Datadog(data map[string]interface{}) {
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
		return
	}

	log.Println(string(out))
}
