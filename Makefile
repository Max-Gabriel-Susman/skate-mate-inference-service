
build:
	docker build -t apollonion-conversation-service .

run: 
	docker run --name apollonion-conversation-service -p 8083:8083 apollonion-conversation-service

connect:
	echo -n “Ground Control For Major Tom” | nc localhost 8083

listen:
	go run cmd/apollonion-server-service/main.go

speak: 
	go run cmd/apollonion-client-service/main.go