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
	"github.com/mchmarny/disco/pkg/target"
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
	query := &types.SimpleQuery{
		OutputPath: reportPath,
		Kind:       types.KindImage,
		Version:    h.version,
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

	if err := h.discoPackages(r.Context(), dir, reportPath); err != nil {
		writeError(w, errors.Wrap(err, "error discovering packages"))
		return
	}

	writeMessage(w, "Done")
}

func (h *Handler) discoLicenses(ctx context.Context, dir, src string) error {
	reportName := fmt.Sprintf("lic-%s.json", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, reportName)
	query := &types.LicenseQuery{
		SimpleQuery: types.SimpleQuery{
			ImageFile:  src,
			OutputPath: reportPath,
			OutputFmt:  types.JSONFormat,
			Kind:       types.KindLicense,
			Version:    h.version,
			TargetRaw:  fmt.Sprintf("bq://%s.disco.licenses", h.projectID),
			Bucket:     h.bucket,
		},
	}

	tar, err := target.ParseImportRequest(&query.SimpleQuery)
	if err != nil {
		return errors.Wrap(err, "error parsing import request")
	}

	if err := disco.DiscoverLicenses(ctx, query, tar); err != nil {
		return errors.Wrap(err, "error executing discover licenses")
	}

	list, err := disco.MeterLicense(ctx, reportPath)
	if err != nil {
		return errors.Wrapf(err, "error metering licenses from: %s", reportPath)
	}

	if err := h.counter.CountAll(ctx, list...); err != nil {
		return errors.Wrapf(err, "error counting licenses metrics: %d", len(list))
	}

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
			TargetRaw:  fmt.Sprintf("bq://%s.disco.vulnerabilities", h.projectID),
			Bucket:     h.bucket,
		},
	}

	tar, err := target.ParseImportRequest(&query.SimpleQuery)
	if err != nil {
		return errors.Wrap(err, "error parsing import request")
	}

	if err := disco.DiscoverVulns(ctx, query, tar); err != nil {
		return errors.Wrap(err, "error executing discover vulnerabilities")
	}

	list, err := disco.MeterVulns(ctx, reportPath)
	if err != nil {
		return errors.Wrapf(err, "error metering vulnerabilities from: %s", reportPath)
	}

	if err := h.counter.CountAll(ctx, list...); err != nil {
		return errors.Wrapf(err, "error counting vulnerability metrics: %d", len(list))
	}

	return nil
}

func (h *Handler) discoPackages(ctx context.Context, dir, src string) error {
	reportName := fmt.Sprintf("pkg-%s.json", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, reportName)
	query := &types.PackageQuery{
		SimpleQuery: types.SimpleQuery{
			ImageFile:  src,
			OutputPath: reportPath,
			OutputFmt:  types.JSONFormat,
			Kind:       types.KindPackage,
			Version:    h.version,
			TargetRaw:  fmt.Sprintf("bq://%s.disco.packages", h.projectID),
			Bucket:     h.bucket,
		},
	}

	tar, err := target.ParseImportRequest(&query.SimpleQuery)
	if err != nil {
		return errors.Wrap(err, "error parsing import request")
	}

	if err := disco.DiscoverPackages(ctx, query, tar); err != nil {
		return errors.Wrap(err, "error executing discover packages")
	}

	list, err := disco.MeterPackage(ctx, reportPath)
	if err != nil {
		return errors.Wrapf(err, "error metering packages from: %s", reportPath)
	}

	if err := h.counter.CountAll(ctx, list...); err != nil {
		return errors.Wrapf(err, "error counting packages metrics: %d", len(list))
	}

	return nil
}
