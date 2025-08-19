package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		ch := getLinesChannel(con)

		for item := range ch {
			fmt.Println(item)
		}
		fmt.Println("Connection to ", con.RemoteAddr(), " closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLineContents := ""
		for {

			b := make([]byte, 8)

			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					return
				}
				return
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				ch <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return ch
}
