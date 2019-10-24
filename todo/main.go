package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"github.com/vkonstantin/wg/todo/controller"
	"github.com/vkonstantin/wg/todo/server/rest"
	"github.com/vkonstantin/wg/todo/storage/memory"
)

var (
	listen = flag.String("listen", ":8080", "Host and port to bind to")
)

func main() {
	flag.Parse()

	storage := memory.NewDefault()
	mainService := controller.NewMainService(storage)
	s := rest.New(*listen, mainService)
	err := s.Run()
	if err != nil {
		log.Printf("error: %s on rest.Run()", err)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Printf("End")
}
