package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/nuuls/log"
)

type Config struct {
	Server   string   `json:"server"`
	Pass     string   `json:"oauth"`
	Nick     string   `json:"username"`
	Channels []string `json:"channels"`
}

func LoadConfig(path string) *Config {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("cannot open config file:", err)
	}
	var c Config
	err = json.Unmarshal(bs, &c)
	if err != nil {
		log.Fatal("error reading config file:", err)
	}
	return &c
}
