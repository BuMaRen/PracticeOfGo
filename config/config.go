package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	MAXCONFIGSIZE = 1024 * 1024
)

type SvrConfig struct {
	DateBaseLocation string  `json:"database"`
	LogFileDir       string  `json:"logpath"`
	SvrHostPort      string  `json:"host:port"`
	ReadTimeout      float64 `json:"readtimeout"`
	WriteTimeout     float64 `json:"writetimeout"`
	DbName           string  `json:"database name"`
	DbContainer      string  `json:"database container"`
}

var config *SvrConfig = nil

func setDefaultConfig() *SvrConfig {
	config.LogFileDir = `/opt/log/wordmemsvrlog`
	config.DateBaseLocation = `mongodb://127.0.0.1:27017`
	config.SvrHostPort = `:20079`
	config.ReadTimeout = time.Duration.Seconds(10)
	config.WriteTimeout = time.Duration.Seconds(10)
	return config
}

func GetConfig() *SvrConfig {
	if config != nil {
		return config
	}
	config = new(SvrConfig)
	fd, err := os.Open("config/svrconfig.json")
	if err != nil {
		fmt.Println("open failed,", err)
		return setDefaultConfig()
	}
	defer fd.Close()
	info, infoerr := fd.Stat()
	buffer := make([]byte, info.Size())
	if l, err := fd.Read(buffer); infoerr != nil || err != nil || l == 0 {
		fmt.Println("read failed,", err)
		return setDefaultConfig()
	}
	if err := json.Unmarshal(buffer, config); err != nil {
		fmt.Println("unmarshal failed,", err)
		return setDefaultConfig()
	}
	return config
}

func (s *SvrConfig) Set() {
	buffer, _ := json.Marshal(config)
	fd, _ := os.OpenFile("svrconfig.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	fd.Write(buffer)
	defer fd.Close()
}
