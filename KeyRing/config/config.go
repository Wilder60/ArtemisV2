package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"Port"`
	} `yaml:"Server"`
	Logger struct {
		Project string `yaml:"Project"`
		Name    string `yaml:"Name"`
	} `yaml:"Logger"`
	Security struct {
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"Security"`
	Database struct {
		Postgres struct {
			Hostname string `yaml:"Hostname"`
			Dbname   string `yaml:"Dbname"`
			User     string `yaml:"User"`
			Password string `yaml:"Password"`
			Port     string `yaml:"Port"`
		} `yaml:"Postgres"`
	} `yaml:"Database"`
}

func ProvideConfig() *Config {
	cfg := &Config{}
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("1 argument required but %d were found", len(os.Args)-1)) // -1 because main file
	}
	cfgPath := os.Args[1]
	file, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		panic(err.Error())
	}

	return cfg
}

var Module = fx.Options(
	fx.Provide(ProvideConfig),
)
