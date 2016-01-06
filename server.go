package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	maxPacketSize = 32 * 1024
)

func startServer(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve addr: %s %v", addr, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to listen for UDP connections: ", err)
	}

	go handleConnection(conn)
	return nil
}

func handleConnection(conn *net.UDPConn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		m, err := parseMetric(scanner.Bytes())
		if err == nil {
			log.Println(err)
		} else {
			fmt.Println(m)
		}
	}
}

func main() {
	_ = startServer("localhost:1119")

	select {}
}
