package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Port string `yaml:"Port"`
	}
	Security struct {
		SecretKey string `yaml:"SecretKey"`
	}
}

var conf *config

// Get will initalize the conf if it is unitalized
func Get() *config {
	if conf != nil {
		conf = new()
	}
	return conf
}

func new() *config {
	path := os.Args[1]
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	c := &config{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(err)
	}
	return c
}
