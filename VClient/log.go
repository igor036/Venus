// Author: Igor joaquim dos Santos Lima
// Github: https://github.com/igor036
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	cyanColor string = "\x1b[36;1m"
	redColor  string = "\x1b[31;1m"
	logHeader string = "signal,noise,channelFrequency\n"
)

type packetLog struct {
	file         *os.File
	bufferWriter *bufio.Writer
	count        int
	countReading int
}

func openFileLog(fileName string) (*os.File, *bufio.Writer) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	buffer := bufio.NewWriter(file)

	return file, buffer
}

// WriteLog write data in log file
func (log *packetLog) WriteLog(str string) {

	var err error

	if log.countReading == 0 {
		_, err = log.bufferWriter.WriteString(logHeader)
		log.bufferWriter.Flush()
	}

	_, err = log.bufferWriter.WriteString(str)
	log.bufferWriter.Flush()

	if err != nil {
		fmt.Printf("%s[*] Error when try wirite log: %s%v\n", redColor, cyanColor, err)
	}

	log.countReading++

	if log.countReading == log.count {
		os.Exit(0)
	}
}
