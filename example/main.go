package main

import (
	"github.com/z0w13/socketgo"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	socket := socketgo.NewServer()
	socket.OnConnect(func(c *socketgo.Conn) error {
		log.Printf("%s connected\n", c.RemoteAddr())
		c.Send("hello", "world")
		return nil
	})
	socket.Handle("hello", func(conn *socketgo.Conn, payload interface{}) error {
		log.Printf("Received 'hello' from client with: %s\n", payload)
		return nil
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, err := ioutil.ReadFile("index.html")
		if err != nil {
			log.Printf("Error reading index.html: %s", err)
		}

		if _, err := w.Write(html); err != nil {
			log.Printf("Error writing response: %s", err)
		}
	})

	mux.HandleFunc("/socketgo.js", func(w http.ResponseWriter, r *http.Request) {
		js, err := ioutil.ReadFile("../socketgo.js")
		if err != nil {
			log.Printf("Error reading socketgo.js: %s", err)
		}

		if _, err := w.Write(js); err != nil {
			log.Printf("Error writing response: %s", err)
		}
	})

	mux.Handle("/ws", socket)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln(err)
	}
}
