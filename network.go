package main

import (
	"os"
	"fmt"
	"net"
)

func ListenForUDPInLoop(port int, recvChannel chan string) {
	println("Creating UDP connection")
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		return
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error receiving UDP connection: %v\n", err)
		return
	}

	println("Created UDP connection")
	for {
		var buf [1024]byte
		dataRead, _, err := conn.ReadFromUDP(buf[:])
		go handleUDPData(buf[:dataRead], err, recvChannel)
	}
}

func handleUDPData(data []byte, err error, recvChannel chan string) {
	fmt.Printf("handleUDPData: %s\n", string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading UDP message: %v\n", err)
	} else {
		recvChannel <- string(data)
	}
}
