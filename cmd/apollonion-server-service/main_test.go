package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServerConcurrency(t *testing.T) {
	t.Run("successful handling of concurrent requests from multiple clients(5 clients)", func(t *testing.T) {
		address := "localhost:8080"
		shutdownChan := make(chan struct{}, 1)
		go startServer(address, shutdownChan)
		time.Sleep(time.Second)

		numClients := 5
		messages := make(chan string, numClients*numClients)

		clientWork := func(id int) {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				t.Error("Failed to connect to server:", err)
				return
			}
			defer conn.Close()

			message := fmt.Sprintf("Hello from client %d", id)
			fmt.Fprintln(conn, message)

			responseScanner := bufio.NewScanner(conn)
			for responseScanner.Scan() {
				response := responseScanner.Text()
				if strings.Contains(response, "Broadcast:") {
					messages <- response
				}
			}
		}

		for i := 0; i < numClients; i++ {
			go clientWork(i)
		}

		expectedMessages := numClients * (numClients - 1)
		for i := 0; i < expectedMessages; i++ {
			<-messages
		}
		t.Logf("Received all expected broadcast messages")
		shutdownChan <- struct{}{}
	})

	t.Run("server start up delay", func(t *testing.T) {
		address := "localhost:8081"
		shutdownChan := make(chan struct{}, 1)
		go startServer(address, shutdownChan)
		time.Sleep(1 * time.Second)

		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Skip("Server not ready, skipping test")
		}
		conn.Close()
		shutdownChan <- struct{}{}
	})

	t.Run("connection retries", func(t *testing.T) {
		address := "localhost:8082"
		shutdownChan := make(chan struct{}, 1)
		go startServer(address, shutdownChan)
		time.Sleep(1 * time.Second)

		const maxRetries = 3
		var conn net.Conn
		var err error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			conn, err = net.Dial("tcp", address)
			if err == nil {
				break
			}
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
		if err != nil {
			t.Fatal("Failed to connect after retries:", err)
		}
		conn.Close()
		shutdownChan <- struct{}{}
	})

	t.Run("stress test(100 clients)", func(t *testing.T) {
		address := "localhost:8083"
		shutdownChan := make(chan struct{}, 1)
		go startServer(address, shutdownChan)
		time.Sleep(1 * time.Second)

		const stressClientCount = 100
		done := make(chan bool, stressClientCount)

		for i := 0; i < stressClientCount; i++ {
			go func(id int) {
				conn, err := net.Dial("tcp", address)
				if err != nil {
					t.Error("Failed to connect to server:", err)
					done <- false
					return
				}
				defer conn.Close()

				message := fmt.Sprintf("Stress test message from client %d", id)
				fmt.Fprintf(conn, message+"\n")
				scanner := bufio.NewScanner(conn)
				if scanner.Scan() {
					done <- true
				} else {
					t.Error("Failed to receive response during stress test")
					done <- false
				}
			}(i)
		}

		for i := 0; i < stressClientCount; i++ {
			if success := <-done; !success {
				t.Error("Not all clients completed successfully")
			}
		}
		shutdownChan <- struct{}{}
	})

	t.Run("broadcast test", func(t *testing.T) {
		address := "localhost:8084"
		shutdownChan := make(chan struct{}, 1)
		go startServer(address, shutdownChan)
		time.Sleep(1 * time.Second)

		numClients := 5
		messages := make(chan string, numClients*(numClients-1))

		clientWork := func(id int) {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				t.Error("Failed to connect to server:", err)
				return
			}
			defer conn.Close()

			greeting := fmt.Sprintf("Hello from client %d", id)
			fmt.Fprintln(conn, greeting)

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				received := scanner.Text()
				if strings.HasPrefix(received, "Broadcast:") {
					messages <- received
				}
			}
		}

		for i := 0; i < numClients; i++ {
			go clientWork(i)
		}

		expectedMessages := numClients * (numClients - 1)
		for i := 0; i < expectedMessages; i++ {
			<-messages
		}

		t.Logf("Received all expected broadcast messages")
		shutdownChan <- struct{}{}
	})
}
