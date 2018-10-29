package main

import ( 
	"bufio"
	"log"
	"os"
	"fmt"
)

const CYAN_COLOR string = "\x1b[36;1m"  
const RED_COLOR  string = "\x1b[31;1m"

type Log struct {

	File 	 	 *os.File
	bufferWriter *bufio.Writer

}

type Config struct {

	DeviceName  string
	LogFile 	*Log

}

func HandleArgs(args []string) Config  {

	if len(args) < 2 { log.Fatal("number of args invalid") }
	
	config := Config {

		DeviceName: "",   
		LogFile:    nil,
	}

	for i := 1; i < len(args); i++ {
		
		arg := args[i]
		i++

		if arg == "-i" { 

			config.DeviceName = args[i]
		
		} else if arg == "-l" { 

			file, err := os.OpenFile(args[i], os.O_WRONLY, 0666)
			if err != nil { log.Fatal(err) }
			buffer := bufio.NewWriter(file)
	
			config.LogFile = &Log{
				File: 		  file,
				bufferWriter: buffer,
			}
		}
	}
	
	if config.DeviceName == "" { log.Fatal("interface name not reported") }

	return config
}


func (log* Log) WriteLog(str string) {

	_, err := log.bufferWriter.WriteString(str)
	log.bufferWriter.Flush()

	if err != nil {
		
		fmt.Printf(
			"%s[*] Error when try wirite log: %s%v\e[0m\n",
			RED_COLOR,
			CYAN_COLOR,
			err,
		)
    }
}