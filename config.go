package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var config Config

type ConfigCollector struct {
	Address     string `json:"address"`
	Port        uint16 `json:"port"`
	Vendor      string `json:"vendor"`
	Version     uint8  `json:"version"`
	Negotiation bool   `json:"negotiation"`
}

type ConfigSession struct {
	Id   byte   `json:"id"`
	Name string `json:"name"`
}

type ConfigExporter struct {
	Address        string          `json:"address"`
	Port           uint16          `json:"port"`
	KeepAlive      uint32          `json:"keep-alive"`
	ConnectTimeout uint32          `json:"connect-timeout"`
	Sessions       []ConfigSession `json:"sessions"`
}

type Config struct {
	Collector ConfigCollector `json:"collector"`
	Exporter  ConfigExporter  `json:"exporter"`
}

func ReadConfig() error {

	jsonFile, err := os.Open("config.json")
	if err != nil {
		return err
	}

	log.Printf("Read config file success!\n")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		log.Printf("Unmarshal failed!\n")
		return err
	}

	return nil
}

func GetServerAddr() string {

	return fmt.Sprintf("%s:%d", config.Exporter.Address, config.Exporter.Port)
}

func GetConnectTimeout() uint32 {
	return config.Exporter.ConnectTimeout
}

func GetConnectParam() (string, uint16, string, uint8, uint32) {
	return config.Collector.Address, config.Collector.Port, config.Collector.Vendor, config.Collector.Version, config.Exporter.KeepAlive
}

func GetSessionList() []byte {
	sessIds := []byte{}
	for _, s := range config.Exporter.Sessions {
		sessIds = append(sessIds, s.Id)
	}

	return sessIds
}
