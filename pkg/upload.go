package pkg

import (
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type Upload struct {
	logger *zerolog.Logger
}

type httpError struct {
	status int
	url    string
}

func CreateUploader(l *zerolog.Logger) *Upload {
	logger := l.With().Str("type", "upload").Logger()
	return &Upload{logger: &logger}
}

func (u *Upload) UploadFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		u.logger.Error().Err(err).Str("url", url).Msg("Failed to connect to URL")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := &httpError{status: resp.StatusCode, url: url}
		u.logger.Error().Err(err).Msg("Non-OK response from URL")
		return nil, err
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error().Err(err).Msg("Error reading image from response body")
		return nil, err
	}

	return imageBytes, nil
}

func (e *httpError) Error() string {
	return "unexpected status code " + http.StatusText(e.status) + " (" + string(rune(e.status)) + ") from " + e.url
}
