package main

import (
	. "boilerplate-go/config"
	_ "boilerplate-go/service"
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		Log.Fatal("Error loading .env file")
	}
	Log.Info("Start Started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-signalChan
	Log.Info("Server Stopped")
}
