package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := "http://localhost:3000"
	for {
		time.Sleep(20 * time.Millisecond)
		go MakeRequest(addr)
	}

}

func MakeRequest(addr string) {
	resp, err := http.Get(addr)
	if err != nil {
		log.Fatalln(err)
	}
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bb))
}
