package pkg

import (
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type Upload struct {
	logger *zerolog.Logger
}

func CreateUploader(l *zerolog.Logger) *Upload {
	logger := l.With().Str("type", "upload").Logger()

	return &Upload{
		logger: &logger,
	}
}

func (u *Upload) UploadFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		u.logger.Error().Err(err).Msg("Failed connect to url: " + err.Error())

		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		u.logger.Error().Err(err).Msg("Received bad response from url: " + err.Error())

		return nil, err
	}

	imageBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		u.logger.Error().Err(err).Msg("Error reading image: " + err.Error())

		return nil, err
	}

	return imageBytes, nil
}
