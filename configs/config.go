package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Port int `yaml:"Port"`
	} `yaml:"Server"`
	Security struct {
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"Security"`
	Database struct {
		SQL struct {
			Project  string `yaml:"Project"`
			Region   string `yaml:"Region"`
			Instance string `yaml:"Instance"`
			Dbname   string `yaml:"Dbname"`
			User     string `yaml:"User"`
			Password string `yaml:"Password"`
		} `yaml:"SQL"`
	} `yaml:"Database"`
}

var conf *config = nil

// Get will initalize the conf if it is unitalized
func Get() *config {
	if conf == nil {
		conf = new()
	}
	return conf
}

func new() *config {
	path := os.Args[1]
	file, err := ioutil.ReadFile(path)
	fmt.Println(string(file))
	if err != nil {
		panic(err)
	}
	c := config{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		panic(err)
	}
	return &c
}
