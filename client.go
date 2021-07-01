package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please, type integer number to send to the server: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nClient IP:\t%s\n", c.LocalAddr())
	fmt.Printf("Server IP:\t%s\n\n", c.RemoteAddr())
	fmt.Println("Waiting for response...")

	fmt.Fprintf(c, text + "\n")
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Message from server: " + message)
	fmt.Println("Connection is closed.")
}
