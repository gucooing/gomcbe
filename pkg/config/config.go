package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	ServerPath  string `json:"ServerPath"`
	CqhttpWsurl string `json:"CqhttpWsurl"`
	QqAdmin     int    `json:"Qqadmin"`
}

var CONF *Config = nil

func GetConfig() *Config {
	return CONF
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	f, err := os.Open(filePath)
	if err != nil {
		return FileNotExist
	}
	defer func() {
		_ = f.Close()
	}()
	c := new(Config)
	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	CONF = c
	return nil
}

var DefaultConfig = &Config{
	ServerPath:  "D:/XM/bds/bedrock_server.exe",
	CqhttpWsurl: "ws://127.0.0.1:8080",
	QqAdmin:     123456789,
}
