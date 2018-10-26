package main

import "log"

type Config struct {

	deviceName string
	trainMode  bool

}

func HandleArgs(args []string) Config {

	if len(args) > 3 || len(args) < 2 { log.Fatal("number of args invalid") }
	
	config := Config {
		deviceName: "",
		trainMode:  false,
	}

	for i := 0; i < len(args); i++ {
		
		arg := args[i]
		i++

		if arg == "-i" { 
			config.deviceName = args[i]
		} else if arg == "-t" { 
			config.trainMode = true
		}
	}
	
	if config.deviceName == "" { log.Fatal("interface name not reported") }

	return config
}