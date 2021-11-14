package internal

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Stream struct {
	Priority int
	URL      string
}

type Config struct {
	Credentials struct {
		Username string `yaml:"username"`
		Secret   string `yaml:"secret"`
	} `yaml:"credentials"`
	Streams []Stream `yaml:"streams"`
}

func LoadConfig() Config {
	f, err := os.Open("/app/config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
