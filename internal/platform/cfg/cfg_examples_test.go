package cfg_test

import (
	"fmt"
	"log"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/cfg"
)

// ExampleInit show how to use the package level func of the config
func ExampleInit() {
	cfg.Init(cfg.MapProvider{
		Map: map[string]string{
			"IP":         "40.23.233.10",
			"PORT":       "4044",
			"INIT_STAMP": time.Date(2009, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(time.UnixDate),
			"FLAG":       "on",
		},
	})

	// To get the ip
	fmt.Println(cfg.MustString("IP"))

	// To get port number
	fmt.Println(cfg.MustInt("PORT"))

	// To get timestamp
	fmt.Println(cfg.MustTime("INIT_STAMP"))

	// To get flag
	fmt.Println(cfg.MustBool("FLAG"))

	// Output:
	// 40.23.233.10
	// 4044
	// 2009-11-10 15:00:00 +0000 UTC
	// true
}

// ExampleNew shows how to create and use a new config
func ExampleNew() {
	c, err := cfg.New(cfg.MapProvider{
		Map: map[string]string{
			"IP":   "80.23.233.10",
			"PORT": "8044",
			"INIT_STAMP": time.Date(2009, time.November,
				10, 23, 0, 0, 0, time.UTC).UTC().Format(time.UnixDate),
			"FLAG": "off",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// To get the ip.
	fmt.Println(c.MustString("IP"))

	// To get the port number.
	fmt.Println(c.MustInt("PORT"))

	// To get the timestamp.
	fmt.Println(c.MustTime("INIT_STAMP"))

	// To get the flag.
	fmt.Println(c.MustBool("FLAG"))

	// Output:
	// 80.23.233.10
	// 8044
	// 2009-11-10 23:00:00 +0000 UTC
	// false
}
