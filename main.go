package main

import (
	"flag"
	apiserver "github.com/mmfc-labs/driving-assistant/pkg/apiserver"
	"github.com/mmfc-labs/driving-assistant/version"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath string
)

func main() {
	log.WithFields(log.Fields{"version": version.Version, "gitRevision": version.GitRevision}).Info("be starting")
	flag.StringVar(&configPath, "config-path", "./config.yaml", "Configuration file path")
	flag.Parse()

	// Run api server
	server := runAPIServer()

	// receive signal exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-quit
	log.Info("signal receive exit ", s)
	server.Stop()
}

// runAPIServer Run apis ,pprof ui ,trace ui
func runAPIServer() *apiserver.APIServer {
	server := apiserver.NewAPIServer(
		apiserver.DefaultOptions().
			WithAddr(":80").
			WithConfigPath(configPath))

	server.Run()
	return server
}
