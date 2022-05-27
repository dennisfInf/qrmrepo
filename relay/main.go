package main

import (
	"github.com/enclaive/relay/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := cmd.Execute(); err != nil {
		log.Fatal().Caller().Err(err).Msg("failed to start")
	}
}
