package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/divideandconquer/go-consul-client/src/config"
)

// envLoader satisfies the Loader interface in go-consul-client
type envLoader struct {
}

// NewEnvLoader creates a Loader that will cache the provided namespace on initialization
// and return data from that cache on Get
func NewEnvLoader() config.Loader {
	return envLoader{}
}

// Import does nothing for the env loader
func (e envLoader) Import(data []byte) error {
	return nil
}

// Initialize does nothing for the env loader
func (e envLoader) Initialize() error {
	return nil
}

// Put does nothing for the env loader
func (e envLoader) Put(key string, value []byte) error {
	return nil
}

// Get fetches the raw config from the environment
func (e envLoader) Get(key string) ([]byte, error) {
	val := os.Getenv(key)
	if val != "" {
		return []byte(val), nil
	}

	return nil, fmt.Errorf("Could not find value for key: %s", key)
}

// MustGetString fetches the config and parses it into a string.  Panics on failure.
func (e envLoader) MustGetString(key string) string {
	b, err := e.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}

	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}

	return s
}

// MustGetBool fetches the config and parses it into a bool.  Panics on failure.
func (e envLoader) MustGetBool(key string) bool {
	b, err := e.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	var ret bool
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}
	return ret
}

// MustGetInt fetches the config and parses it into an int.  Panics on failure.
func (e envLoader) MustGetInt(key string) int {
	b, err := e.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}

	var ret int
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}
	return ret
}

// MustGetDuration fetches the config and parses it into a duration.  Panics on failure.
func (e envLoader) MustGetDuration(key string) time.Duration {
	s := e.MustGetString(key)
	ret, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("Could not parse config (%s) into a duration: %v", key, err))
	}
	return ret
}
