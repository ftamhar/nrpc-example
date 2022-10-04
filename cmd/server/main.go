package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example_nrpc/proto/hello"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type server struct{}

func (s *server) Upload(ctx context.Context, req *hello.UploadRequest) (*hello.UploadResponse, error) {
	name := uuid.NewString() + ".jpeg"
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	b := req.GetData()
	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res := hello.UploadResponse{
		Name: name,
	}

	return &res, nil
}

func (s *server) Greeting(ctx context.Context, req *hello.GreetingRequest) (resp *hello.GreetingResponse, err error) {
	resp = &hello.GreetingResponse{
		Fullname: req.Firstname + " " + req.Lastname,
	}
	return
}

func main() {
	natsURL := nats.DefaultURL
	// Connect to the NATS server.
	// nc, err := nats.Connect(natsURL, nats.Timeout(5*time.Second))
	nc, err := nats.Connect(natsURL, nats.Timeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Our server implementation.
	s := &server{}

	// The NATS handler from the helloworld.nrpc.proto file.
	h := hello.NewHelloServicesHandler(context.TODO(), nc, s)
	sub, err := nc.QueueSubscribe(h.Subject(), h.Subject(), h.Handler) // subject = HelloServices.>
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	fmt.Println("server is running, ^C quits.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	close(c)
}
