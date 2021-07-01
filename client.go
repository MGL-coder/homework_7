package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// reading number from os.Stdin to send
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please, type integer number that will be send to the server: ")
	textToSend, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// connecting to server
	c, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	// sending text to the server
	fmt.Fprintf(c, textToSend + "\n")

	fmt.Printf("\nClient IP:\t%s\n", c.LocalAddr())
	fmt.Printf("Server IP:\t%s\n\n", c.RemoteAddr())
	fmt.Println("Waiting for response...")

	// getting response
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Message from server: " + message)
	fmt.Println("Connection is closed.")
}
