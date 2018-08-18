package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// Container contains the information about the container
type Container struct {
	ID   string
	Port string
}

// StartMongo runs a mongo container to execute commnads
func StartMongo() (*Container, error) {
	cmd := exec.Command("docker", "run", "-P", "-d", "mongo")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("starting container: %v", err)
	}

	id := out.String()[:12]
	log.Println("DB ContatinerID:", id)
	cmd.Wait()

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("inspect container: %v", err)
	}

	cmd.Wait()

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP27017 []struct {
					HostPort string `json:"HostPort"`
				} `json:"27017/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		return nil, fmt.Errorf("decoding json: %v", err)
	}
	c := Container{
		ID:   id,
		Port: doc[0].NetworkSettings.Ports.TCP27017[0].HostPort,
	}

	log.Println("DB Port: ", c.Port)

	return &c, nil
}

// StopMongo stops and removes specified container
func StopMongo(c *Container) error {
	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		return err
	}
	log.Println("Stopped: ", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		return err
	}

	log.Println("Removed:", c.ID)
	return nil
}
