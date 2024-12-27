package instrumentation

import (
	"backend-layout/internal/config"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(conf *config.Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(conf.App.LogLevel)

	if err != nil {
		log.Warn().Err(err).Msgf("Invalid log level '%s', defaulting to INFO", conf.App.LogLevel)
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var stdOut io.Writer = os.Stdout
	if conf.App.Env == "development" {
		stdOut = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	logger := zerolog.New(stdOut).With().Timestamp().Caller().Logger()
	log.Logger = logger
	return logger
}
