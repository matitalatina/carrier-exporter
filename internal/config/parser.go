package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Carriers struct {
		Wind struct {
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
			LineID     string `yaml:"lineId"`
			ContractID string `yaml:"contractId"`
		} `yaml:"wind"`
		Tim struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Phone    string `yaml:"phone"`
		} `yaml:"tim"`
		Vodafone struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Sims     []struct {
				Phone    string   `yaml:"phone"`
				Counters []string `yaml:"counters"`
			} `yaml:"sims"`
		} `yaml:"vodafone"`
	} `yaml:"carriers"`
}

func ParseFile(path string) (Config, error) {
	config := Config{}

	absPath, err := filepath.Abs(path)

	if err != nil {
		return config, err
	}

	data, err := ioutil.ReadFile(absPath)

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}
