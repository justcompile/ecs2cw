package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type account struct {
	ID     string `json:"id"`
	Region string `json:"region"`
}

type Config struct {
	Accounts []*account `json:"accounts"`
}

func NewConfig(filePath string) (*Config, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cfg *Config

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
