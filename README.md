# Nuance Retrieval Service 

## Overview 

The Nuance Retrieval Service is exposes a RAG application over a tcp server

# EKS

The scripts will create resources in: us-east-1

# How to run it
- export AWS_PROFILE=user1 # The name of the profile you want to use
- ./create-vpc-stack.sh
- ./create-eks-stack.sh

# How to clean up
- ./delete-eks-stack.sh
- ./delete-vpc-stack.sh

# Tips

## Update Kubeconfig
aws eks update-kubeconfig --region us-east-1 --name my-eks-cluster

## Install an ingress controller
https://kubernetes.github.io/ingress-nginx/deploy/#quick-start

## Verify ingress controller
kubectl get service ingress-nginx-controller --namespace=ingress-nginx

[05/01/2024 notes]

wrap goroutines in an anonymous function always

mutex uneeded here because:

we're synchronously boradcasting message

so we'll have to wait for that message to be sent to every connection

parallel/async?

We need to actually write it back to the client sort of like an event for the client to receive

net.Conn has a writer maybe? it has a reader and the writer, uses bytes()

instead of mutex, in broadcaster() use a sync.Waitgroup and then fanout the inner for loop(faster)

just parallize inner for sake of simplicity(outer should be irl tho)
    - waitgroup instead of mutex
    - cap concurrency to n goroutines(use a buffered channel as a semaphore)
    - 

handleConnection:
    use conn.Read() function instead of scanner.Scan

we're gonna need a goroutine per connection

can be optimized

potential race condition: you potentially send to a connection after it's been closed before it's been removed from connection slice

global lock defeats purpose of concurrency

don't lock all connections, lock the slice when appending and deleting 

you want to have individual locks on each connection 

in the context of concurrency most data races come from two process writing to the same data at the same time

might not need to here

rework locking/unlocking strategy for sake of performance

so we can easily test and mock this we should better organize this into packages and then create an interface for the client that defines read and write so we can moq it up and make this generally better 

shoot for production readiness as a chat server

CONDENSED:


wrap goroutines in an anonymous function always



handleConnection:
    use conn.Read() function instead of scanner.Scan

we're gonna need a goroutine per connection


so we can easily test and mock this we should better organize this into packages and then create an interface for the client that defines read and write so we can moq it up and make this generally better 

[05/10/2024]

we need a way to gracefully shutdown
- prerfably a channel with a length of 1 that we can return from start server and write to to tell it to end the for loop

- use select for the graceful shutdown implementation to handle different cases

* actual work with concurrency in the marketplace domain
- good opportunity 

also put interfaces by their consumers as opposed to the structs their resembling

[05/17/2024]

* create a wrapper around conns so that chat features can be implemented 

* implement feedback for the sender on if the recipient received the message 

* broadcast a message sent to all other chatters in the server

[05/31/2024]


[07/17/2024]

implemente horizontal skaling w/ k8s(2 or more)

get client to connect to one of the servers and then the responding server broadcasts to all clients

put a load balancing application in front of the two servers

eventing service for handling the messages 

can speed up implementation by leveraging k8s from docker desktop 

test coverage as well 

NewServer() should start the listener and can then return itself plus the error so you can still check the error in main

pass the config to new server as well so you can parameterize port of service

we should have a distinct client package

have a start function(will be your handle client thing) and then call that with a goroutine in main, and then the start function will have your main for loop in it

go s.HandleClient(conn) should be in Start or Begin function

simplify server

implement concurrent broadcast to clients, needs to handle for arbitray scaling 

fence off open API stuff

proper commenting, documentation, and testing 
