/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
*/
package main

import (
	"os"
	"net"
	"log"
	"fmt"
  "time"
	"os/exec"
  "github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

var (
	config Config
	conn   net.Conn
)

const DOT_11_POBRE_REQUEST  layers.Dot11Type = 0x10

type Dot_11_Info struct {

	SrcAddress 			 net.HardwareAddr
	DstAddress 			 net.HardwareAddr
	Type						 string
	Bssid						 []byte
	ChannelFrequency uint16
	Signal					 int8
	Noise						 int8

}

func MonitorMode() {
	
	cmds := [3]string { 
		"sudo ifconfig wlp1s0 down", 
		"sudo iwconfig wlp1s0 mode monitor",
		"sudo ifconfig wlp1s0 up",
	}
	
	for _, cmd := range cmds {

		_, err := exec.Command("sh","-c",cmd).Output()

		if (err != nil) { log.Fatal(err)  }
	}
}


func CreateDot_11_Info(packet  gopacket.Packet) *Dot_11_Info {

	dot11 := packet.Layer(layers.LayerTypeDot11)
	radio := packet.Layer(layers.LayerTypeRadioTap)
	
	if nil != dot11 && radio != nil { 
		
		dot11, _ := dot11.(*layers.Dot11) 
		radio, _ := radio.(*layers.RadioTap)

		var dotType string 
		var bssid  []byte 
	
		if dot11.Type == DOT_11_POBRE_REQUEST {

			dotType = "Pobre Request"

		} else {

			dotType = "Pobre Response"
			
			dot11info := packet.Layer(layers.LayerTypeDot11InformationElement)
			
			if nil != dot11info {
				dot11info, _ := dot11info.(*layers.Dot11InformationElement)
				if dot11info.ID == layers.Dot11InformationElementIDSSID { bssid = dot11info.Info }
			}
		}

		//false packet
		if radio.DBMAntennaSignal >= 0x0  { return nil }

		return &Dot_11_Info { 

			SrcAddress: 		dot11.Address2,
			DstAddress:			dot11.Address1,	
			ChannelFrequency: 	uint16(radio.ChannelFrequency),
			Type:				dotType,
			Bssid:				bssid,
			Signal:				radio.DBMAntennaSignal,
			Noise: 				radio.DBMAntennaNoise,

		}
	}

	return nil

}


func HandlerPkt(packet  gopacket.Packet) {

	dot_11_Info := CreateDot_11_Info(packet)

	if dot_11_Info != nil {

		fmt.Printf("Type: %s\n", dot_11_Info.Type)
		fmt.Printf("BSSID: %q\n", dot_11_Info.Bssid)
		fmt.Printf("SRC Address: %v\n", dot_11_Info.SrcAddress)
		fmt.Printf("DST Address: %v\n", dot_11_Info.DstAddress)
		fmt.Printf("Frequency: %d\n", dot_11_Info.ChannelFrequency)
		fmt.Printf("Signal: %ddbm\n", dot_11_Info.Signal)
		fmt.Printf("Noise: %ddbm\n", dot_11_Info.Noise)
		fmt.Printf("\n\n")

		data := fmt.Sprintf(
			"\t[ %d, %d, %d ],\n",
			dot_11_Info.Signal,
			dot_11_Info.Noise,
			dot_11_Info.ChannelFrequency,
		)

		if config.CanWriteLog(dot_11_Info.SrcAddress) { 
			config.LogFile.WriteLog(data) 
		} else { 
			conn.Write([]byte (data )) 
		}
	}
}

func Start() {

	MonitorMode()

	config = HandleConfig()

	if config.LogMode {
		defer config.LogFile.File.Close()
	} else {
		conn = Connection(config.ServerAddress)
		defer conn.Close()
	}

    handle, err := pcap.OpenLive(config.DeviceName, 1024, false, 30 * time.Second)
	
    if err != nil { log.Fatal(err) }
	defer handle.Close()

	err = handle.SetBPFFilter("type mgt subtype probe-req or type mgt subtype probe-resp")
	if err != nil { log.Fatal(err) }

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet   := range packetSource.Packets() { HandlerPkt(packet) }

}

func main() {

	if len(os.Args) > 1 {
		RunArg(os.Args[1])
	} else {
		Start()
	}

}