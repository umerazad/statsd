package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	maxPacketSize = 32 * 1024
)

var stats = newAdminStats()
var agg = newAggregator()

func startMetricsServer(addr string) {
	fmt.Println("Starting metrics server at: ", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve addr: %s %v", addr, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to listen for UDP connections: ", err)
	}

	go handleConnection(conn)
}

func handleConnection(conn *net.UDPConn) {
	buf := make([]byte, maxPacketSize)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("Failed to read from conn: %v", err)
		}

		for _, r := range strings.Split(string(buf[:n]), "\n") {
			r = strings.TrimSpace(r)
			for _, r := range strings.Split(string(buf[:n]), "\n") {
				if len(r) > 0 {
					m, err := parseMetric([]byte(r))
					if err != nil {
						stats.BadRecords++
						log.Println(err)
					} else {
						stats.TotalRecords++
						switch m.Type {
						case COUNTER:
							stats.CounterRecords++
							agg.processCounter(m)
						case GUAGE:
							stats.GuageRecords++
							agg.processGuage(m)
						case TIMING:
							stats.TimingRecords++
							agg.processTiming(m)
						case SET:
							stats.SetRecords++
							agg.processSet(m)
						default:
							// Just for debugging
							log.Fatalf("Unrecognized metric: %v", m)
						}
					}
				}
			}
		}
	}
}

func startAdminServer(addr string) {
	fmt.Println("Starting admin server at: ", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start admin server: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection on admin interface: %v", err)
		}

		stats.TotalAdminConnections++
		go handleAdminConn(conn)
	}
}

func handleAdminConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 128)
	for {
		_, err := c.Read(buf)
		if err == nil {
			fmt.Print("Admin channel received: ", string(buf))
		}
	}
}

func startFlusher(seconds time.Duration) {
	fmt.Println("Starting flusher with period:", seconds)
	for {
		select {
		case <-time.After(seconds):
			agg.writeCounters(os.Stdout)
			agg.writeGuages(os.Stdout)
			agg.writeTimings(os.Stdout)
		}
	}
}

func main() {
	go startMetricsServer("localhost:1119")
	go startAdminServer("localhost:2000")
	go startFlusher(15 * time.Second)
	select {}
}
