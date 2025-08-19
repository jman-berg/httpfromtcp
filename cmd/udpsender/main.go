package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const udpAddressString = "localhost:42069"

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp", udpAddressString)
	if err != nil {
		log.Fatalf("Error resolving UDP-address %s", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		log.Fatalf("Error setting up UPD connection on address: %s", udpAddressString)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Print("Error reading input:", err.Error())
			return
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatalf("Error writing to udp-connection: %s", err.Error())
		}

	}

}
