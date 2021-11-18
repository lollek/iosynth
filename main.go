package main

import (
	"flag"
	"log"
)

func main() {
	var udpPort int

	flag.IntVar(&udpPort, "p", 49161, "UDP port to use")
	flag.Parse()

	recvChannel := make(chan []byte)
	go ListenForUDPInLoop(udpPort, recvChannel)

	if err := InitSoundServer(); err != nil {
		log.Fatalf("Sound server failed to start: %v", err)
	}

	for {
		data, ok := <-recvChannel
		if !ok {
			log.Fatal("UDP Channel has closed")
		} else {
			HandleData(data)
		}
	}
}
