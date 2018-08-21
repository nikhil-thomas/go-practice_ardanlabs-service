package publisher

import (
	"encoding/json"
	"log"
)

// Console handles metrics for direct display on console
func Console(data map[string]interface{}) {
	out, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return
	}
	log.Println(string(out))
}
