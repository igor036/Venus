/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
*/
package main

import (
	"fmt"
	"bufio" 
	"os"
	"log"
	"net"
	"bytes"
	"io/ioutil"
	"os/exec"
	"encoding/json"
)

const ( 
	JSON_FILE_NAME 	     string  = "_config.json" 
	DEFAULT_CONTENT_JSON string = 
`{	
	"deviceName":   "wlp1s0",
	"logFileName":  "log.txt",
	"logAddress":   "00:00:00:00:00:00",
	"logCount":      2,
	"logMode":       true,
	"serverAddress": "127.0.0.1:180"
}`
)

type ConfigJson struct {

	DeviceName    string `json:"deviceName"`
	LogFileName   string `json:"logFileName"`
	LogAddress    string `json:"logAddress"`
	ServerAddress string `json:"serverAddress"`
	LogCount      int    `json:"logCount"`
	LogMode       bool   `json:"logMode"`

}

type Config struct {

	DeviceName    string
	LogMode       bool
	LogAddress    net.HardwareAddr
	ServerAddress string
	LogFile       *Log

}

func (config* Config) CanWriteLog(addr net.HardwareAddr) bool {
	
	if addr == nil || config.LogAddress == nil { return false }

	return bytes.Equal(addr, config.LogAddress)
}

func ReadConfigDotJson() ConfigJson {

	var configJson ConfigJson
	 
  data, err := ioutil.ReadFile(JSON_FILE_NAME)
  if err != nil { log.Fatal(err) }
  
  err = json.Unmarshal(data, &configJson)
	if err != nil { log.Fatal(err) }

	return configJson

}

func RunArg(arg string) {

	argFunctions := map[string]func() {
		
		"--open-config": func() {
			_, err := exec.Command("sh","-c",fmt.Sprintf("sudo gedit %s",JSON_FILE_NAME)).Output()
			if err != nil { log.Fatal(err)  }
		},

		"--restore-config": func() {
			
			file, err := os.OpenFile(JSON_FILE_NAME, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil { log.Fatal(err) }

			buffer := bufio.NewWriter(file)
			_, err = buffer.WriteString(DEFAULT_CONTENT_JSON)
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

func HandleConfig() Config  {
	
	configJson := ReadConfigDotJson()
	
	//log address
	hwAddr, err := net.ParseMAC(configJson.LogAddress)
	if err != nil { log.Fatal(err) }

	//create log and buffer write file
	file, err := os.OpenFile(configJson.LogFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil { log.Fatal(err) }

	buffer := bufio.NewWriter(file)

	config := Config {
		DeviceName:    configJson.DeviceName, 
		LogMode:	   configJson.LogMode,
		LogAddress:	   hwAddr,
		ServerAddress: configJson.ServerAddress,
		LogFile: 	   &Log {
			File: 		  file,
			bufferWriter: buffer,
			Count:		  configJson.LogCount,
			CountReading: 0,
		},
	}


	return config
}