package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "timeout")

	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := host + ":" + port

	c := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := c.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer func(client TelnetClient) {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}(c)

	log.Printf("...Connected to %s \n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		<-ctx.Done()
		cancel()
	}()

	go func() {
		if err := c.Send(); err != nil {
			log.Fatalln(err)
		}

		log.Print("...EOF")
		cancel()
	}()

	go func() {
		err := c.Receive()
		if err != nil {
			log.Fatalln(err)
		}

		log.Print("...Connection was closed by peer")
		cancel()
	}()
}
