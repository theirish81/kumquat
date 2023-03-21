package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v1 "github.com/theirish81/kumquat/api/v1"
	"os"
)

func main() {
	if os.Getenv("MODE") == "dev" {
		log.Level(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Info().Msg("Kumquat starting")
	api, _ := v1.NewAPI()
	api.Run()
}
