package disco

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/disco/pkg/metric"
	"github.com/pkg/errors"
)

// result is the handler response type.
type result struct {
	Status  string `json:"status"`
	Image   string `json:"image,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewHandler creates new handler instance.
func NewHandler(bucket string, counter metric.Counter) (*Handler, error) {
	if bucket == "" {
		return nil, errors.New("bucket name not set")
	}

	if counter == nil {
		return nil, errors.New("counter service not set")
	}

	h := &Handler{
		bucket:  bucket,
		counter: counter,
	}

	return h, nil
}

// Handler is the handler type.
type Handler struct {
	bucket  string
	counter metric.Counter
}

// HandlerDefault is the default handler.
func (h *Handler) HandlerDefault(w http.ResponseWriter, r *http.Request) {
	writeMessage(w, "nothing to see here, try /disco")
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println(err)
	writeContent(w, result{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  err.Error(),
	})
}

func writeMessage(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	log.Println(msg)
	writeContent(w, result{
		Status:  http.StatusText(http.StatusOK),
		Message: msg,
	})
}

func writeContent(w http.ResponseWriter, content any) {
	if err := json.NewEncoder(w).Encode(content); err != nil {
		log.Printf("error encoding: %v - %v", content, err)
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
