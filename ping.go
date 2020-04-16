package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)
var doc = `Usage:
1. go build ping.go
2. sudo ./ping [-t timeout] [-c count] hostname/ip

Examples:
# ping google continuously
sudo ./ping www.google.com

# ping google 10 times
sudo ./ping -c 10 www.google.com

# ping IP address 127.123.0.168 continuously
sudo ./ping 127.123.0.168

# ping google with time limit in milliseconds
sudo ./ping -t 500 google.com 
`

func getMessage(t icmp.Type) icmp.Message {
	return icmp.Message {
		Type: t,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Seq: 1,
			Data: []byte("WRITTEN BY: NHAN LE . VISIT ME @ nhantle.com"),
		},
	}
}

func getPacket() *icmp.PacketConn {
	packet, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return packet
}

func ping(dest string, counter int, timeOut int) bool {
	//timeout boolean
	isTimedOut := false

	//encode "echo request" message
	message := getMessage(ipv4.ICMPTypeEcho)
	packetBytes, err := message.Marshal(nil)
	if err != nil {
		fmt.Println("There has been error, please try again")
		return true
	}

	//get packet from the caller
	packet := getPacket()

	//destination ip address
	destAddr, err := net.ResolveIPAddr("ip4:icmp", dest)
	if err != nil {
		fmt.Println("There has been error, please try again")
		return true
	}

	beginTime := time.Now()

	//send the packet to the destinationAddr IP
	_, err = packet.WriteTo(packetBytes, destAddr)
	if err != nil {
		fmt.Println(err)
		return true
	}

	// get the reply from the listener
	messageBytes := make([]byte, 1500)
	packet.ReadFrom(messageBytes)

	endTime := time.Since(beginTime)

	// logging info to the terminal
	replyMessage, err := icmp.ParseMessage(1, messageBytes)
	if err != nil {
		fmt.Println(err)
	}

	// Print out stats
	fmt.Printf("%d bytes from %s icmp_seq=%d [Time: %s]\n",replyMessage.Body.Len(0), destAddr, counter,endTime)

	// check time out
	if endTime.Milliseconds() - int64(timeOut) > 0{
		fmt.Println("Operation Timeout")
		isTimedOut = true
	}
	packet.Close()
	return isTimedOut
}

func main() {
	timeOut := flag.Int("t", 1000 , "")
	count := flag.Int("c", int(^uint(0) >> 1), "")
	flag.Parse()

	if len(flag.Args()) == 0{
		fmt.Print(doc)
	}

	dest := flag.Arg(0)

	//ping forever with counter
	for i:= 0; i < *count; i++ {
		if ping(dest, i, *timeOut){
			return
		}
		time.Sleep(time.Second)
	}
}

