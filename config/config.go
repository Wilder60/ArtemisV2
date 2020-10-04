package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"Port"`
	}
	Security struct {
		SecretKey string `yaml:"SecretKey"`
	}
}

var conf *Config

// Get will initalize the conf if it is unitalized
func Get() *Config {
	if conf != nil {
		conf = new()
	}
	return conf
}

func new() *Config {
	path := os.Args[1]
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	c := &Config{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(err)
	}
	return c
}
