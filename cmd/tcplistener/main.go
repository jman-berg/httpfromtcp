package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jman-berg/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening for tcp traffic: %s", err.Error())
	}
	defer l.Close()

	fmt.Println("Listening for TCP traffic on port:", port)

	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		fmt.Println("Accepted connection from: ", con.RemoteAddr())

		parsedRequest, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatalf("Error parsing request: %s", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", parsedRequest.RequestLine.Method)
		fmt.Println("- Target:", parsedRequest.RequestLine.RequestTarget)
		fmt.Println("- Version:", parsedRequest.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, _ := range parsedRequest.Headers {
			fmt.Printf("- %s: %v\n", k, parsedRequest.Headers[k])
		}
		fmt.Println("Body:")
		fmt.Println(string(parsedRequest.Body))

		fmt.Println("Connection to ", con.RemoteAddr(), " closed")
	}

}
