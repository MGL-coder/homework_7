package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Client IP:\t%s\n", c.LocalAddr())
	fmt.Printf("Server IP:\t%s\n\n", c.RemoteAddr())

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Text to send: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(c, text + "\n")
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Message from server: " + message)
	fmt.Println("Connection is closed.")
}
