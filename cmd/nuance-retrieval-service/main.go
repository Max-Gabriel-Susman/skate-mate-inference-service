package main

import (
	"fmt"
	"net"
	"os"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Max-Gabriel-Susman/nuance-retrieval-service/internal/server"
)

type Config struct {
	EmbeddingModel string
	InfernecModel  string
	OpenaiKey      string
	PineKey        string
}

func main() {
	// ctx := context.Background()

	cfg := Config{
		// EmbeddingModel: "text-embedding-ada-002",
		// InfernecModel:  "gpt-3.5-turbo",
		OpenaiKey: os.Getenv("OPENAI_API_KEY"),
		// PineKey:        os.Getenv("PINECONE_API_KEY"),
	}

	// pineKey := os.Getenv("PINECONE_API_KEY")

	openAPIClient := openai.NewClient(cfg.OpenaiKey)

	// pc, err := pinecone.NewClient(pinecone.NewClientParams{
	// 	ApiKey: pineKey,
	// })

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// idxs, err := pc.ListIndexes(ctx)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// for _, index := range idxs {
	// 	fmt.Println(index)
	// }

	// idx, err := pc.Index(idxs[0].Host)
	// defer idx.Close()

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// res, err := idx.DescribeIndexStats(ctx)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// fmt.Println(res)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := server.NewServer(listener, openAPIClient)
	defer s.Listener.Close()

	go s.HandleConnections()

	fmt.Println("Server listening on port 8080")
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.HandleClient(conn)
	}
}
