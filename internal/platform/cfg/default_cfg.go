package cfg

import (
	"net/url"
	"time"
)

// c is the default Config used by Innit and the package level ffuncs
// like String, MustString, and SetString
var c Config

// Init populates the package default Config and should be called only once.
// A Provider must be supplied which will return a map of key/value pairs to be loaded
func Init(p Provider) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, err := p.Provide()
	if err != nil {
		return err
	}
	c.m = m

	return nil
}

// Log returns a string to help with a package's default config
// excludes any values whose key contains the string PASS
func Log() string {
	return c.Log()
}

// String works on default config
func String(key string) (string, error) {
	return c.String(key)
}

// MustString works on default config
func MustString(key string) string {
	return c.MustString(key)
}

// SetString works on default config
func SetString(key, value string) {
	c.SetString(key, value)
}

// Int works on default config
func Int(key string) (int, error) {
	return c.Int(key)
}

// MustInt works on default config
func MustInt(key string) int {
	return c.MustInt(key)
}

// SetInt works on default config
func SetInt(key string, value int) {
	c.SetInt(key, value)
}

// Time works on default config
func Time(key string) (time.Time, error) {
	return c.Time(key)
}

// MustTime works on default config
func MustTime(key string) time.Time {
	return c.MustTime(key)
}

// SetTime works on default config
func SetTime(key string, value time.Time) {
	c.SetTime(key, value)
}

// Bool works on default config
func Bool(key string) (bool, error) {
	return c.Bool(key)
}

// MustBool works on default config
func MustBool(key string) bool {
	return c.MustBool(key)
}

// SetBool works on default config
func SetBool(key string, value bool) {
	c.SetBool(key, value)
}

// URL works on default config
func URL(key string) (*url.URL, error) {
	return c.URL(key)
}

// MustURL works on default config
func MustURL(key string) *url.URL {
	return c.MustURL(key)
}

// SetURL works on default config
func SetURL(key string, value *url.URL) {
	c.SetURL(key, value)
}

// Duration works on default config
func Duration(key string) (time.Duration, error) {
	return c.Duration(key)
}

// MustDuration works on default config
func MustDuration(key string) time.Duration {
	return c.MustDuration(key)
}

// SetDuration works on default config
func SetDuration(key string, value time.Duration) {
	c.SetDuration(key, value)
}
