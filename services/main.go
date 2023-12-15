package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var count int64 = 0
var Port string

func main() {
	port := os.Getenv("SERVER_PORT")
	Port = port
	log.Printf("New service started on port: %s", port)
	app := fiber.New()
	app.Get("/", handlefunc)
	timer := os.Getenv("CLOSE_TIMEOUT")

	var time_ int
	var err error
	if timer != "" {
		time_, err = strconv.Atoi(timer)
		if err != nil {
			log.Fatal(err)
		}
	}

	go appStop(app, time_)
	go countPrinter(port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(time_) * time.Second)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}

func countPrinter(port string) {
	for {
		time.Sleep(5 * time.Second)
		log.Printf("Port %s | Conns %d", port, atomic.LoadInt64(&count))
	}
}

func handlefunc(c *fiber.Ctx) error {
	atomic.AddInt64(&count, 1)
	time.Sleep(500 * time.Millisecond)
	atomic.AddInt64(&count, -1)
	return c.SendString(fmt.Sprintf("Success query %s", Port))
}

func appStop(a *fiber.App, time_ int) {
	if time_ == 0 {
		return
	}
	time.Sleep(time.Duration(time_) * time.Second)
	a.Shutdown()
}
