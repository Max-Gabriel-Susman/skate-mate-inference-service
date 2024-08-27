package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	connMutex     sync.Mutex
	connections   []net.Conn
	broadcastChan = make(chan string)
)

func broadcaster(shutdownChan chan struct{}) {
	for {
		select {
		case msg := <-broadcastChan:
			connMutex.Lock()
			for _, conn := range connections {
				fmt.Fprintf(conn, "Broadcast: %s\n", msg)
			}
			connMutex.Unlock()
		case <-shutdownChan:
			return
		}
	}
}

func handleConnection(conn net.Conn, shutdownChan chan struct{}) {
	defer func() {
		conn.Close()
		connMutex.Lock()
		for i, c := range connections {
			if c == conn {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		connMutex.Unlock()
	}()

	connMutex.Lock()
	connections = append(connections, conn)
	connMutex.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case broadcastChan <- scanner.Text():
			fmt.Println("Message received:", scanner.Text())
		case <-shutdownChan:
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading from connection: %s\n", err)
	}
}

func startServer(address string, shutdownChan chan struct{}) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("Server is listening on", address)

	go broadcaster(shutdownChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-shutdownChan:
				return nil
			default:
				fmt.Fprintf(os.Stderr, "error accepting connection: %s\n", err)
				continue
			}
		}
		go handleConnection(conn, shutdownChan)
	}
}

func main() {
	shutdownChan := make(chan struct{}, 1)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("Shutdown signal received, shutting down...")
		shutdownChan <- struct{}{}

		connMutex.Lock()
		for _, conn := range connections {
			conn.Close()
		}
		connections = nil
		connMutex.Unlock()
	}()

	if err := startServer(":8080", shutdownChan); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %s\n", err)
		os.Exit(1)
	}
}
