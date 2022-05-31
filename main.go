package main

import (
	"github.com/mmfc-labs/driving-assistant/apiserver"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Run api server
	apiServer := runAPIServer()

	// receive signal exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-quit
	log.Info("signal receive exit ", s)
	apiServer.Stop()
}

// runAPIServer Run apis ,pprof ui ,trace ui
func runAPIServer() *apiserver.APIServer {
	apiServer := apiserver.NewAPIServer(
		apiserver.DefaultOptions().
			WithAddr(":8080"))

	apiServer.Run()
	return apiServer
}
