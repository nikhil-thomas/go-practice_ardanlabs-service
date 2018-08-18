package cfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// EnvProvider provides configuration from the environment
// All keys will be made uppercase
type EnvProvider struct {
	Namespace string
}

// Provide implements the Provider interface
func (ep EnvProvider) Provide() (map[string]string, error) {

	config := map[string]string{}

	// Get all available environment variables
	envs := os.Environ()

	if len(envs) == 0 {
		return nil, errors.New("No environment variables found")
	}

	uNameSpace := fmt.Sprintf("%s_", strings.ToUpper(ep.Namespace))

	for _, val := range envs {
		if !strings.HasPrefix(val, uNameSpace) {
			continue
		}
		idx := strings.Index(val, "=")
		config[strings.ToUpper(strings.TrimPrefix(val[0:idx], uNameSpace))] = val[idx+1:]
	}
	if len(config) == 0 {
		return nil, fmt.Errorf("Namespace %q was not found", ep.Namespace)
	}
	return config, nil
}
