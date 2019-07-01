// Author: Igor joaquim dos Santos Lima
// Github: https://github.com/igor036
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
)

const (
	configJSONFileName       string = "resources/config.json"
	defaultContentConfigJSON string = `{	
		"deviceName":   "wlp1s0",
		"logFileName":  "log.txt",
		"logAddress":   "00:00:00:00:00:00",
		"logCount":      2,
		"logMode":       true,
		"serverAddress": "127.0.0.1:180"
	}`
)

type configJSON struct {
	deviceName    string `json:"deviceName"`
	logFileName   string `json:"logFileName"`
	logAddress    string `json:"logAddress"`
	serverAddress string `json:"serverAddress"`
	logCount      int    `json:"logCount"`
	logMode       bool   `json:"logMode"`
}

// Config properties of app
type Config struct {
	DeviceName    string
	LogMode       bool
	LogAddress    net.HardwareAddr
	ServerAddress string
	LogFile       *Log
}

// IsLogAddress verify if a Address
// is the Address of device using for log
func (config *Config) IsLogAddress(addr net.HardwareAddr) bool {

	if addr == nil || config.LogAddress == nil {
		return false
	}

	return bytes.Equal(addr, config.LogAddress)
}

func readConfigDotJSON() configJSON {

	var configJSON configJSON

	data, err := ioutil.ReadFile(configJSONFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &configJSON)
	if err != nil {
		log.Fatal(err)
	}

	return configJSON

}

func runArg(arg string) {

	argFunctions := map[string]func(){

		"--open-config": func() {
			_, err := exec.Command("sh", "-c", fmt.Sprintf("sudo gedit %s", configJSONFileName)).Output()
			if err != nil {
				log.Fatal(err)
			}
		},

		"--restore-config": func() {

			file, err := os.OpenFile(configJSONFileName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}

			buffer := bufio.NewWriter(file)
			_, err = buffer.WriteString(defaultContentConfigJSON)
			buffer.Flush()

		},
	}

	fun := argFunctions[arg]

	if fun != nil {
		fun()
	} else {
		log.Fatal("Invalid arg")
	}
}

func handleConfig() Config {

	configJSON := readConfigDotJSON()
	file, buffer := openFileLog(configJSON.logFileName)

	hwAddr, err := net.ParseMAC(configJSON.logAddress)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{
		DeviceName:    configJSON.deviceName,
		LogMode:       configJSON.logMode,
		LogAddress:    hwAddr,
		ServerAddress: configJSON.serverAddress,
		LogFile: &Log{
			file:         file,
			bufferWriter: buffer,
			count:        configJSON.logCount,
			countReading: 0,
		},
	}

	return config
}
