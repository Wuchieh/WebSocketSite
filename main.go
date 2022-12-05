package main

import (
	"github.com/Wuchieh/WebSocketSite/Server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sc := make(chan os.Signal, 1)
	go func() {
		err := Server.Server()
		if err != nil {
			log.Println(err)
			sc <- os.Kill
			return
		}
	}()

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
