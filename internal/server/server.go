package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Max-Gabriel-Susman/nuance-retrieval-service/internal/message"
)

type Client struct {
	channel chan string
	name    string
}

type ServerProvider interface {
	HandleClient(conn net.Conn)
	HandleConnections()
}

type Server struct {
	Listener     net.Listener
	Clients      map[Client]bool
	Broadcast    chan string
	Register     chan Client
	Unregister   chan Client
	Mutex        *sync.Mutex
	OpenAIClient *openai.Client
}

func NewServer(listener net.Listener, openAIClient *openai.Client) Server {
	return Server{
		Listener:     listener,
		Clients:      make(map[Client]bool),
		Broadcast:    make(chan string),
		Register:     make(chan Client),
		Unregister:   make(chan Client),
		Mutex:        &sync.Mutex{},
		OpenAIClient: openAIClient,
	}

}

func (s Server) HandleConnections() {
	for {
		select {
		case text := <-s.Broadcast:
			s.Mutex.Lock()
			msg := message.NewMessage(text)
			resp, err := msg.RespondToMessage(s.OpenAIClient)
			if err != nil {
				log.Printf("ChatCompletion error: %v\n", err)
				continue
			}
			fmt.Println(resp.Choices[0].Message.Content)
			for client := range s.Clients {
				select {
				case client.channel <- resp.Choices[0].Message.Content:
				default:
					close(client.channel)
					delete(s.Clients, client)
				}
			}
			s.Mutex.Unlock()
		case client := <-s.Register:
			s.Mutex.Lock()
			s.Clients[client] = true
			s.Mutex.Unlock()
		case client := <-s.Unregister:
			s.Mutex.Lock()
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.channel)
			}
			s.Mutex.Unlock()
		}
	}
}

func (s Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	channel := make(chan string)
	client := Client{channel: channel, name: conn.RemoteAddr().String()}

	s.Register <- client

	go func() {
		for text := range channel {
			fmt.Fprintln(conn, text)
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		s.Broadcast <- fmt.Sprintf("%s: %s", client.name, text)
	}

	s.Unregister <- client
	fmt.Printf("Client %s disconnected\n", client.name)
}
