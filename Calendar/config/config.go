package config

import (
	"io/ioutil"
	"os"

	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

// Config is the exported struct mimicing
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
		Firebase struct {
			ServiceAccount string `yaml:"ServiceAccount"`
		} `yaml:"Firebase"`
	} `yaml:"Database"`
}

// CreateConfig will create a new Config file
func CreateConfig() *Config {
	if len(os.Args) != 2 {
		panic("")
	}

	cfgfile := os.Args[1]
	data, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		panic("Failed to read configuration file")
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic("")
	}
	return cfg
}

var ConfigModule = fx.Option(
	fx.Provide(CreateConfig),
)
