package main

import (
	"github.com/prb-releases/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Release Server Starting...")
	config := server.Default()
	s := server.NewServer(config)

	s.Run()
}
