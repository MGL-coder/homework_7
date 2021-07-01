package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type Server struct {
	listener net.Listener
	quit chan struct{}
	wg sync.WaitGroup
	sm chan struct{}	// semaphore to restrict the max number of connections that server can handle at the same time
}

func newServer(port string) *Server {
	fmt.Println("Starting new server...")
	s := &Server{
		quit: make(chan struct{}),
	}
	l, err := net.Listen("tcp4", ":" + port)
	if err != nil {
		log.Fatalf("network listen initialization error: %s", err)
	}
	s.listener = l

	fmt.Println("Please enter the max number of connections that server can handle at the same time:")
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("cannot read from standard input: %s",err)
	}

	n, err := strconv.Atoi(strings.TrimSpace(temp))
	if err != nil {
		log.Fatalf("incorrect input: integer expected: %s", err)
	}
	s.sm = make(chan struct{}, n)

	s.wg.Add(1)
	go s.serve()
	return s

}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		c, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.sm <- struct{}{}
				s.handleConnection(c)
				<- s.sm
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	defer func() {fmt.Printf("Connection with %s is closed\n", c.RemoteAddr().String())}()

	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	netDataNumber, err := strconv.Atoi(strings.TrimSpace(netData))
	if err != nil {
		_, err := c.Write([]byte("incorrect input: integer expected\n"))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	result := strconv.Itoa(netDataNumber * netDataNumber) + "\n"

	_, err = c.Write([]byte(result))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func main() {
	s := newServer("8081")
	fmt.Println("Server is running...")
	fmt.Println("Press \"enter\" to shutdown the server at any time.")

	quit := make(chan os.Signal, 1)

	// graceful shutdown caused by input
	go func() {
		fmt.Scanln()
		close(quit)
	}()

	// graceful shutdown caused by syscall
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	fmt.Println("Server shutting down...")
	s.stop()
	fmt.Println("Server stopped.")
}
