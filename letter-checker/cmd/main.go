package main

import (
	"letter-checker/config"
	"letter-checker/requester"

	"fmt"
	"letter-checker/api"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	time.Sleep(30 * time.Second)

	start := time.Now()
	cfg := config.Load("config")

	configLogging(cfg)
	printLogHeader(cfg)
	cfg.Print()

	repo := requester.CassandraChecker{}
	repo.Init(cfg)
	
		
	defer repo.Close()

	service := requester.NewCheckerService(&repo)
	r := api.ConfigureRouter(service)

	log.Info().Str("port", cfg.Port.Value).Int64("startTimeMs", time.Since(start).Milliseconds()).Msg("listening")
	log.Fatal().Err(http.ListenAndServe(":"+cfg.Port.Value, r))
}

func printLogHeader(cfg *config.Config) {
	if cfg.Log.Structured.Value {
		log.Info().Str("application", cfg.AppName.Value).
			Str("revision", cfg.Revision.Value).
			Str("version", cfg.AppVersion.Value).
			Str("sha1ver", cfg.Sha1Version.Value).
			Str("build-time", cfg.BuildTime.Value).
			Str("profile", cfg.Profile.Value).
			Str("config-source", cfg.Config.Source.Value).
			Send()
	} else {
		f := figure.NewFigure(cfg.AppName.Value, "", true)
		f.Print()

		log.Info().Msg("=============================================")
		log.Info().Msg(fmt.Sprintf("      Revision: %s", cfg.Revision.Value))
		log.Info().Msg(fmt.Sprintf("       Profile: %s", cfg.Profile.Value))
		log.Info().Msg(fmt.Sprintf("   Tag Version: %s", cfg.AppVersion.Value))
		log.Info().Msg(fmt.Sprintf("  Sha1 Version: %s", cfg.Sha1Version.Value))
		log.Info().Msg(fmt.Sprintf("    Build Time: %s", cfg.BuildTime.Value))
		log.Info().Msg("=============================================")
	}
}



func configLogging(cfg *config.Config) {
	log.Info().Msg("configuring logging...")

	if !cfg.Log.Structured.Value {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	level, err := zerolog.ParseLevel(cfg.Log.Level.Value)
	if err != nil {
		log.Warn().Str("loglevel", cfg.Log.Level.Value).Err(err).Msg("defaulting to info")
		level = zerolog.InfoLevel
	}
	log.Info().Str("loglevel", level.String()).Msg("setting log level")
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}