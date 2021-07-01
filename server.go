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
	"time"
)

type Server struct {
	listener net.Listener	// to listen to specific address
	quit chan struct{}		// to signal the server to stop
	wg sync.WaitGroup		// to wait completion of processing of all requests
	sm chan struct{}		// semaphore to restrict the max number of requests that can be processed simultaneously
}

const responseTime = 10 // time in seconds needed to respond to one request

func newServer(port string) *Server {
	fmt.Println("Starting new server...")

	// initializing quit
	s := &Server{
		quit: make(chan struct{}),
	}

	// initializing listener to specific port
	l, err := net.Listen("tcp4", ":" + port)
	if err != nil {
		log.Fatalf("network listen initialization error: %s", err)
	}
	s.listener = l

	// initializing semaphore
	fmt.Println("Please enter the max number of requests that can be processed simultaneously:")
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("cannot read from standard input: %s",err)
	}

	n, err := strconv.Atoi(strings.TrimSpace(temp))
	if err != nil {
		log.Fatalf("incorrect input: integer expected: %s", err)
	}
	if n < 1 {
		log.Fatalf("incorrect input: the number cannot be less than 1.")
	}
	s.sm = make(chan struct{}, n)

	// starting the server
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
	time.Sleep(time.Second * responseTime)		// waiting...

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
	fmt.Println("\nServer is running...")
	fmt.Printf("Press \"enter\" to shutdown the server at any time.\n\n")

	quit := make(chan os.Signal, 1)

	// graceful shutdown caused by syscall
	// Note: server may not be able to process all current request within shutdown timeout by syscall
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// graceful shutdown caused by input
	go func() {
		fmt.Scanln()
		close(quit)
	}()

	<-quit

	fmt.Printf("Server shutting down...\n\n")
	s.stop()
	fmt.Printf("\nServer stopped.\n")
}
