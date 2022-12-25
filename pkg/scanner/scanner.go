package scanner

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/mchmarny/disco/pkg/scanner/trivy"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	LicenseScanner ScannerType = iota
	VulnerabilityScanner
)

var (
	ScanLicense  MakeLicenseCmd = trivy.MakeLicenseCmd
	ParseLicense LicenseParser  = trivy.ParseLicenses

	ScanVulnerability  MakeVulnerabilityCmd = trivy.MakeVulnerabilityCmd
	ParseVulnerability VulnerabilityParser  = trivy.ParseVulnerabilities
)

type ScannerType int64

func (s ScannerType) String() string {
	switch s {
	case LicenseScanner:
		return "license"
	case VulnerabilityScanner:
		return "vulnerability"
	default:
		return "unknown"
	}
}

// MakeLicenseCmd is an interface for license scanners.
type MakeLicenseCmd func(digest, path string) *exec.Cmd

// LicenseParser is an interface for license parsers.
type LicenseParser func(image, path string, filter types.ItemFilter) (*types.LicenseReport, error)

// MakeVulnerabilityCmd is an interface for vulnerability scanners.
type MakeVulnerabilityCmd func(digest, path string) *exec.Cmd

// VulnerabilityParser is an interface for vulnerability parsers.
type VulnerabilityParser func(image, path string, filter types.ItemFilter) (*types.VulnerabilityReport, error)

// GetLicenses returns licenses for the given image.
func GetLicenses(digest, path string, filter types.ItemFilter) (*types.LicenseReport, error) {
	cmd := ScanLicense(digest, path)
	if err := runCmd(cmd, path); err != nil {
		return nil, errors.Wrap(err, "error running license scanning command")
	}

	report, err := ParseLicense(digest, path, filter)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing licenses")
	}

	return report, nil
}

// GetVulnerabilities returns vulnerabilities for the given image.
func GetVulnerabilities(digest, path string, filter types.ItemFilter) (*types.VulnerabilityReport, error) {
	cmd := ScanVulnerability(digest, path)
	if err := runCmd(cmd, path); err != nil {
		return nil, errors.Wrap(err, "error running vulnerability scanner command")
	}

	report, err := ParseVulnerability(digest, path, filter)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing vulnerabilities")
	}

	return report, nil
}

func runCmd(cmd *exec.Cmd, path string) error {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msgf("error running command: %s (out: %s, err: %s)",
			cmd, outb.String(), errb.String())
		return errors.Wrap(err, "error running command")
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Error().Err(err).Msgf("out: %s, err: %s", outb.String(), errb.String())
		return errors.Wrapf(err, "expected report file not found: %s", path)
	}

	return nil
}
