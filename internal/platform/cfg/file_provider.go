package cfg

import (
	"bufio"
	"os"
	"strings"
)

// FileProvider describes a file based loader to load configuration
type FileProvider struct {
	FileName string
}

// Provide implements the Provider interface
func (fp FileProvider) Provide() (map[string]string, error) {

	var config = make(map[string]string)

	file, err := os.Open(fp.FileName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < 3 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		index := strings.Index(line, "=")

		if index <= 0 {
			continue
		}

		if index == len(line)-1 {
			continue
		}

		config[line[:index]] = line[index+1:]
	}

	return config, nil
}
