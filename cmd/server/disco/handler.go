package disco

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mchmarny/disco/pkg/metric"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// result is the handler response type.
type result struct {
	Status  string `json:"status"`
	Image   string `json:"image,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewHandler creates new handler instance.
func NewHandler(version, bucket string, counter metric.Counter) (*Handler, error) {
	if bucket == "" || version == "" {
		return nil, errors.New("bucket name or version not set")
	}

	if counter == nil {
		return nil, errors.New("recorder service not set")
	}

	h := &Handler{
		version: version,
		bucket:  bucket,
		counter: counter,
	}

	return h, nil
}

// Handler is the handler type.
type Handler struct {
	version string
	bucket  string
	counter metric.Counter
}

// HandlerDefault is the default handler.
func (h *Handler) HandlerDefault(w http.ResponseWriter, r *http.Request) {
	writeMessage(w, "nothing to see here, try /disco")
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Error().Err(err).Msg("error")
	writeContent(w, result{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  err.Error(),
	})
}

func writeMessage(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	log.Info().Msg(msg)
	writeContent(w, result{
		Status:  http.StatusText(http.StatusOK),
		Message: msg,
	})
}

func writeContent(w http.ResponseWriter, content any) {
	if err := json.NewEncoder(w).Encode(content); err != nil {
		log.Error().Msgf("error encoding: %v - %v", content, err)
	}
}

func makeFolder(sha string) (string, error) {
	p := fmt.Sprintf("./%s", sha)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(p, os.ModePerm); err != nil {
			return "", errors.Wrap(err, "error creating folder")
		}
	}
	return p, nil
}
