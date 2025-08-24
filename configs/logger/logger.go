package logger

import (
	"github.com/rs/zerolog"
	"os"
	"runtime/debug"
	"time"
)

var (
	log zerolog.Logger
)

func init() {
	buildInfo, _ := debug.ReadBuildInfo()

	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

}

func Info(message string) {
	log.Info().Msg(message)
}

func Warn(message string) {
	log.Warn().Msg(message)
}

func Error(message string, err error) {
	if err != nil {
		log.Error().Err(err).Msg(message)
	}
	log.Error().Msg(message)
}

func Fatal(message string, err error) {
	log.Fatal().Err(err).Msg(message)
}
