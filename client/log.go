/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
*/
package main

import ( 
	"log"
	"bufio"
	"os"
	"fmt"
)

const CYAN_COLOR string = "\x1b[36;1m"  
const RED_COLOR  string = "\x1b[31;1m"

type Log struct {

	File 	 	 *os.File
	bufferWriter *bufio.Writer
	Count		 int
	CountReading int

}

func OpenFileLog(fileName string) (*os.File, *bufio.Writer ) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil { log.Fatal(err) }

	buffer := bufio.NewWriter(file)

	return file, buffer;
}

func (log* Log) WriteLog(str string) {

	var err error

	if log.CountReading == 0 {

		_, err = log.bufferWriter.WriteString("[\n")
		log.bufferWriter.Flush()

	}

	_, err = log.bufferWriter.WriteString(str)
	log.bufferWriter.Flush()

	if err != nil {
		
		fmt.Printf("%s[*] Error when try wirite log: %s%v\n",RED_COLOR,CYAN_COLOR,err)
	}
	
	log.CountReading ++

	if log.CountReading == log.Count {

		_, err = log.bufferWriter.WriteString("]")
		log.bufferWriter.Flush()
		os.Exit(0)

	}
}