package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := connectToServer("localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	sendAndReceiveMessages(conn)
}

func connectToServer(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

func sendAndReceiveMessages(conn net.Conn) {
	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send, type 'exit' to stop:")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "exit" {
			break
		}
		_, err := fmt.Fprintf(conn, text+"\n")
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

func readMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Print("Received from server: ", message)
	}
}
