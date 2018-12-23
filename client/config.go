package main

import ( 
	"bufio"
	"log"
	"os"
	"fmt"
	"net"
	"bytes"
)

const CYAN_COLOR string = "\x1b[36;1m"  
const RED_COLOR  string = "\x1b[31;1m"

type Log struct {

	File 	 	 *os.File
	bufferWriter *bufio.Writer

}

type Config struct {

	DeviceName  string
	Address		net.HardwareAddr
	LogFile 	*Log

}

func (config* Config) EqualAdress(addr net.HardwareAddr) bool {

	if addr == nil || config.Address == nil {
		return false
	}

	return bytes.Equal(addr, config.Address)

}

func (config* Config) CanWrite(addr net.HardwareAddr) bool {
	return config.LogFile != nil && config.EqualAdress(addr)
}

func (log* Log) WriteLog(str string) {

	_, err := log.bufferWriter.WriteString(str)
	log.bufferWriter.Flush()

	if err != nil {
		
		fmt.Printf("%s[*] Error when try wirite log: %s%v\n",RED_COLOR,CYAN_COLOR,err)
    }
}

func HandleArgs(args []string) Config  {

	if len(args) < 2 { log.Fatal("number of args invalid") }
	
	config := Config {

		DeviceName: "", 
		Address:	nil,
		LogFile:    nil,
	}

	for i := 1; i < len(args); i++ {
		
		arg := args[i]
		i++

		if arg == "-i" { 

			config.DeviceName = args[i]
		
		} else if arg == "-a" {

			hwAddr, err := net.ParseMAC(args[i])
			if err != nil { log.Fatal(err) }

			config.Address = hwAddr

		} else if arg == "-l" { 

			file, err := os.OpenFile(args[i], os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil { log.Fatal(err) }

			buffer := bufio.NewWriter(file)
	
			config.LogFile = &Log{
				File: 		  file,
				bufferWriter: buffer,
			}
		}
	}
	
	if config.DeviceName == "" { log.Fatal("interface name not reported") }
	if config.LogFile 	 == nil { log.Fatal("log filel name not reported") }
	if config.Address 	 == nil { log.Fatal("adress not reported") }
	
	return config
}
