package main

import (
	"log"
	"flag"
)

func main() {
	var udpPort int

	flag.IntVar(&udpPort, "p", 49161, "UDP port to use")
	flag.Parse()

	recvChannel := make(chan []byte)
	go ListenForUDPInLoop(udpPort, recvChannel)

	for {
		data, ok := <-recvChannel
		if !ok {
			log.Fatal("UDP Channel has closed")
		}
		log.Printf("Data received: %v\n", data)
	}
}
