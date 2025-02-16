package instrumentation

import (
	"backend-layout/internal/config"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitializeLogger(conf config.AppConfig) func() {
	level, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(level)

	var stdOut io.Writer = os.Stdout
	// if conf.Type == "text" {
	// 	stdOut = zerolog.ConsoleWriter{Out: os.Stdout}
	// }
	writers := []io.Writer{stdOut}
	var runLogFile *os.File
	// if conf.LogFileEnabled {
	runLogFile, err = os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}

	writers = append(writers, runLogFile)
	// }

	zerolog.TimeFieldFormat = time.RFC3339Nano

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return func() {
		if runLogFile != nil {
			runLogFile.Close()
		}
	}
}
