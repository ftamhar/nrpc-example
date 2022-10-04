package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"example_nrpc/proto/hello"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}
	// Connect to the NATS server.
	nc, err := nats.Connect(natsURL, nats.Timeout(5*time.Second), nats.SyncQueueLen(1))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	cli := hello.NewHelloServicesClient(nc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := cli.Greeting(r.Context(), &hello.GreetingRequest{
			Firstname: "Rahmat",
			Lastname:  "Fathoni",
		})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("coba"))
			return
		}
		w.Write([]byte(res.Fullname))
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("golang.jpg")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("coba"))
			return
		}
		res, err := cli.Upload(r.Context(), &hello.UploadRequest{
			Data: b,
		})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("coba"))
			return
		}

		w.Write([]byte(res.GetName()))
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
