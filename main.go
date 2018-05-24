package main

import (
	. "boilerplate-go/config"
	_ "boilerplate-go/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	Log.Info("Start Started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-signalChan
	Log.Info("Server Stopped")
}
