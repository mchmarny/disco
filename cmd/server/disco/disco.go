package disco

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/object"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DiscoHandler is the HTTP handler for disco service.
func (h *Handler) DiscoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Debug().Msg("preparing discovery...")

	requestContextID := uuid.NewString()
	dir, err := makeFolder(requestContextID)
	if err != nil {
		writeError(w, errors.Wrapf(err, "error creating context from: %s", requestContextID))
		return
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Error().Msgf("error deleting context: %s", dir)
		}
	}()

	// get image IDs
	imageReport := fmt.Sprintf("img-%s.txt", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, imageReport)
	query := &types.ImagesQuery{
		URIOnly: true,
		SimpleQuery: types.SimpleQuery{
			OutputPath: reportPath,
			Kind:       types.KindImage,
			Version:    h.version,
		},
	}

	if err := disco.DiscoverImages(r.Context(), query); err != nil {
		writeError(w, errors.Wrap(err, "error discovering images"))
		return
	}

	if err := h.discoVulns(r.Context(), dir, reportPath); err != nil {
		writeError(w, errors.Wrap(err, "error discovering vulnerabilities"))
		return
	}

	if err := h.discoLicenses(r.Context(), dir, reportPath); err != nil {
		writeError(w, errors.Wrap(err, "error discovering licenses"))
		return
	}

	writeMessage(w, "Done")
}

func (h *Handler) discoLicenses(ctx context.Context, dir, src string) error {
	reportName := fmt.Sprintf("lic-%s.json", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, reportName)
	query := &types.SimpleQuery{
		ImageFile:  src,
		OutputPath: reportPath,
		OutputFmt:  types.JSONFormat,
		Kind:       types.KindVulnerability,
		Version:    h.version,
	}

	if err := disco.DiscoverLicenses(ctx, query); err != nil {
		return errors.Wrap(err, "error executing discover licenses")
	}

	if err := object.Save(ctx, h.bucket, reportName, reportPath); err != nil {
		return errors.Wrapf(err, "error writing content to: %s/%s",
			h.bucket, reportName)
	}

	list, err := disco.MeterLicense(ctx, reportPath)
	if err != nil {
		return errors.Wrapf(err, "error metering licenses from: %s", reportPath)
	}

	if err := h.counter.CountAll(ctx, list...); err != nil {
		return errors.Wrapf(err, "error counting licenses metrics: %d", len(list))
	}

	log.Info().Msgf("license report saved to: gs://%s/%s", h.bucket, reportName)

	return nil
}

func (h *Handler) discoVulns(ctx context.Context, dir, src string) error {
	reportName := fmt.Sprintf("vul-%s.json", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, reportName)
	query := &types.VulnsQuery{
		SimpleQuery: types.SimpleQuery{
			ImageFile:  src,
			OutputPath: reportPath,
			OutputFmt:  types.JSONFormat,
			Kind:       types.KindVulnerability,
			Version:    h.version,
		},
	}

	if err := disco.DiscoverVulns(ctx, query); err != nil {
		return errors.Wrap(err, "error executing discover vulnerabilities")
	}

	if err := object.Save(ctx, h.bucket, reportName, reportPath); err != nil {
		return errors.Wrapf(err, "error writing content to: %s/%s",
			h.bucket, reportName)
	}

	list, err := disco.MeterVulns(ctx, reportPath)
	if err != nil {
		return errors.Wrapf(err, "error metering vulnerabilities from: %s", reportPath)
	}

	if err := h.counter.CountAll(ctx, list...); err != nil {
		return errors.Wrapf(err, "error counting vulnerability metrics: %d", len(list))
	}

	log.Info().Msgf("vulnerability report saved to: gs://%s/%s", h.bucket, reportName)

	return nil
}
