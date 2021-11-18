package main

import (
	"log"
	"fmt"
	"net"
)

func ListenForUDPInLoop(port int, recvChannel chan []byte) {
	log.Printf("Creating UDP connection\n")
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error resolving UDP address: %v\n", err)
		return
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("Error receiving UDP connection: %v\n", err)
		return
	}

	log.Printf("Created UDP connection\n")
	for {
		var buf [1024]byte
		dataRead, _, err := conn.ReadFromUDP(buf[:])
		go handleUDPData(buf[:dataRead], err, recvChannel)
	}
}

func handleUDPData(data []byte, err error, recvChannel chan []byte) {
	if err != nil {
		log.Printf("Error reading UDP message: %v\n", err)
	} else {
		recvChannel <- data
	}
}
