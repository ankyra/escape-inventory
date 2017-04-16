package main

import (
	"encoding/json"
	"io/ioutil"
    "os"
)

type Config struct {
    Secret string `json:"secret"`
    configPath string `json:"-"`
}

func NewConfig(file string) *Config {
    return &Config{
        Secret: "waddup",
        configPath: file,
    }

}

func loadConfig(file string) (*Config, error) {
	var config Config

    if !PathExists(file) {
        return NewConfig(file), nil
    }

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
