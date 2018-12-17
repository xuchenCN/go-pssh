package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xuchenCN/go-pssh/cmd"
)

func main() {
	command := cmd.NewPsshCommand()

	if err := command.Execute(); err != nil {
		log.Error(err)
	}
}
