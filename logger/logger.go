package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func New(level string) (logger zerolog.Logger, err error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	loglevel, err := zerolog.ParseLevel(level)

	if err != nil {
		return
	}

	zerolog.SetGlobalLevel(loglevel)

	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	zerolog.DefaultContextLogger = &logger

	return
}
