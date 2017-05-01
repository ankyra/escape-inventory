package config

import (
	"encoding/json"
    "strings"
	"io/ioutil"
	"os"
    "fmt"
    "gopkg.in/yaml.v2"
)

type DatabaseSettings struct {
    Path string `json:"path" yaml:"path"`
}

type StorageSettings struct {
    Path string `json:"path" yaml:"path"`
    Bucket string `json:"bucket" yaml:"bucket"`
    Credentials map[string]string `json:"credentials" yaml:"credentials"`
}

type Config struct {
    Database string `json:"database" yaml:"database"`
    DatabaseSettings DatabaseSettings `json:"database_settings" yaml:"database_settings"`
    StorageBackend string `json:"storage_backend" yaml:"storage_backend"`
    StorageSettings StorageSettings `json:"storage_settings" yaml:"storage_settings"`
}

func NewConfig() *Config {
	return &Config{
        Database: "sqlite",
        DatabaseSettings: DatabaseSettings {
            Path: "/var/lib/escape/registry.db",
        },
        StorageBackend: "local",
        StorageSettings: StorageSettings {
            Path: "/var/lib/escape/releases",
        },
	}

}

func LoadConfig(file string) (*Config, error) {
	var config Config

	if !PathExists(file) {
		return nil, fmt.Errorf("The referenced configuration file '%s' could not be found", file)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
        return nil, fmt.Errorf("Error reading configuration file '%s': %s", file, err.Error())
	}

    if strings.HasSuffix(file, ".yaml") {

        if err = yaml.Unmarshal(b, &config); err != nil {
            return nil, fmt.Errorf("Could not parse YAML in configuration file '%s': %s", file, err.Error())
        }
    } else {
        if err = json.Unmarshal(b, &config); err != nil {
            return nil, fmt.Errorf("Could not parse JSON in configuration file '%s': %s", file, err.Error())
        }
    }

	return &config, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
