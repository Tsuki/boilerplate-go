package main

import (
	log "github.com/sirupsen/logrus"
	_ "boilerplate-go/config"
	_ "boilerplate-go/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Info("Start Started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-signalChan
	log.Info("Server Stopped")
}
