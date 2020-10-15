package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	// FilePath is the location to configuration file.
	FilePath = "./splable.yml"
)

// File is used to model our configuration file.
type File struct {
	Environment string `yaml:"environment"`
	Hostname    string `yaml:"hostname"`
	Token       string `yaml:"token"`
}

// GetConf reads a yml file and loads it into struct.
func (c *File) GetConf() *File {
	fileName, _ := filepath.Abs(FilePath)
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading config file: #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	return c
}
