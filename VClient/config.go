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
	DeviceName    string `json:"deviceName"`
	LogFileName   string `json:"logFileName"`
	LogAddress    string `json:"logAddress"`
	ServerAddress string `json:"serverAddress"`
	LogCount      int    `json:"logCount"`
	LogMode       bool   `json:"logMode"`
}

type configProperties struct {
	deviceName    string
	logMode       bool
	logAddress    net.HardwareAddr
	serverAddress string
	log           *packetLog
}

func (config *configProperties) isLogAddress(addr net.HardwareAddr) bool {

	if addr == nil || config.logAddress == nil {
		return false
	}

	return bytes.Equal(addr, config.logAddress)
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
				fmt.Println("Erro ao abrir arquivo de configuração")
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

func handleConfig() configProperties {

	configJSON := readConfigDotJSON()
	file, buffer := openFileLog(configJSON.LogFileName)

	hwAddr, err := net.ParseMAC(configJSON.LogAddress)
	if err != nil {
		log.Fatal(err)
	}

	config := configProperties{
		deviceName:    configJSON.DeviceName,
		logMode:       configJSON.LogMode,
		logAddress:    hwAddr,
		serverAddress: configJSON.ServerAddress,
		log: &packetLog{
			file:         file,
			bufferWriter: buffer,
			count:        configJSON.LogCount,
			countReading: 0,
		},
	}

	return config
}
