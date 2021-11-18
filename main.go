package main

import (
	"fmt"
	"flag"
)

func main() {
	var udpPort int

	flag.IntVar(&udpPort, "p", 49161, "UDP port to use")
	flag.Parse()

	recvChannel := make(chan string)
	go ListenForUDPInLoop(udpPort, recvChannel)

	for {
		data, ok := <-recvChannel
		if !ok {
			println("Channel closed")
			return
		}
		fmt.Printf("Data received: %v\n", data)
	}
}
