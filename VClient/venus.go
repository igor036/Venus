// Author: Igor joaquim dos Santos Lima
// Github: https://github.com/igor036
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	config configProperties
	conn   net.Conn
)

type dot11Info struct {
	srcAddress       net.HardwareAddr
	dstAddress       net.HardwareAddr
	channelFrequency uint16
	signal           int8
	noise            int8
}

func monitorMode() {

	cmds := [3]string{
		fmt.Sprintf("sudo ifconfig %s down", config.deviceName),
		fmt.Sprintf("sudo iwconfig %s mode monitor", config.deviceName),
		fmt.Sprintf("sudo ifconfig %s up", config.deviceName),
	}

	for _, cmd := range cmds {
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createDot11Info(packet gopacket.Packet) *dot11Info {

	dot11 := packet.Layer(layers.LayerTypeDot11)
	radio := packet.Layer(layers.LayerTypeRadioTap)

	if nil != dot11 && radio != nil {

		dot11, _ := dot11.(*layers.Dot11)
		radio, _ := radio.(*layers.RadioTap)

		//false packet
		if radio.DBMAntennaSignal == 0x0 {
			return nil
		}

		return &dot11Info{
			srcAddress:       dot11.Address2,
			dstAddress:       dot11.Address1,
			channelFrequency: uint16(radio.ChannelFrequency),
			signal:           radio.DBMAntennaSignal,
			noise:            radio.DBMAntennaNoise,
		}
	}

	return nil

}

func handlerPkt(packet gopacket.Packet) {

	dot11Info := createDot11Info(packet)

	if dot11Info != nil && isHandlerPkt(dot11Info) {

		fmt.Printf("SRC Address: %v\n", dot11Info.srcAddress)
		fmt.Printf("DST Address: %v\n", dot11Info.dstAddress)
		fmt.Printf("Frequency: %d\n", dot11Info.channelFrequency)
		fmt.Printf("Signal: %ddbm\n", dot11Info.signal)
		fmt.Printf("Noise: %ddbm\n", dot11Info.noise)
		fmt.Printf("\n\n")

		data := fmt.Sprintf(
			"%d, %d, %d, %d\n",
			dot11Info.signal,
			dot11Info.noise,
			dot11Info.channelFrequency,
			config.distance,
		)

		if config.logMode && config.isLogAddress(dot11Info.srcAddress) {
			config.log.WriteLog(data)
		} else if !config.logMode {
			conn.Write([]byte(data))
		}
	}
}

func isHandlerPkt(dot11Info *dot11Info) bool {
	return config.logMode && config.isLogAddress(dot11Info.srcAddress) || !config.logMode
}

func start() {

	config = handleConfig()
	monitorMode()

	if config.logMode {
		defer config.log.file.Close()
	} else {
		conn = Connection(config.serverAddress)
		defer conn.Close()
	}

	handle, err := pcap.OpenLive(config.deviceName, 1024, false, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	err = handle.SetBPFFilter("type mgt subtype probe-req")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlerPkt(packet)
	}
}

func main() {
	if len(os.Args) > 1 {
		runArg(os.Args[1])
	} else {
		start()
	}
}
