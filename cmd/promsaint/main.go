package main

import (
	"flag"

	_ "github.com/engineyard/promsaint/logging"
	"github.com/engineyard/promsaint/server"
	log "github.com/Sirupsen/logrus"
)

var (
	Version   string
	BuildTime string
)

func main() {
	flag.Parse()

	log.WithFields(log.Fields{
		"version": Version,
		"build":   BuildTime,
	}).Warn("Starting Promsaint!")
	s := server.NewPromsaint()
	s.Start()
}
