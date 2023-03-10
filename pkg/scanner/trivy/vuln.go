package trivy

import (
	"fmt"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type vulnerabilityReport struct {
	Image   string `json:"ArtifactName"` //nolint:tagliatelle
	Results []struct {
		Vulnerabilities []struct {
			ID             string `json:"VulnerabilityID"`  //nolint:tagliatelle
			URL            string `json:"PrimaryURL"`       //nolint:tagliatelle
			Package        string `json:"PkgName"`          //nolint:tagliatelle
			PackageVersion string `json:"InstalledVersion"` //nolint:tagliatelle
			Title          string `json:"Title"`            //nolint:tagliatelle
			Description    string `json:"Description"`      //nolint:tagliatelle
			Severity       string `json:"Severity"`         //nolint:tagliatelle
			Updated        string `json:"LastModifiedDate"` //nolint:tagliatelle
		} `json:"Vulnerabilities"` //nolint:tagliatelle
	} `json:"Results"` //nolint:tagliatelle
}

func ParseVulnerabilities(path string, filter types.ItemFilter) (*types.VulnerabilityReport, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	log.Debug().Msgf("parsing vulnerabilities from %s using filter: %v", path, filter)

	var report vulnerabilityReport
	if err := types.UnmarshalFromFile(path, &report); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file: %s", path)
	}

	list := make([]*types.Vulnerability, 0)
	m := make(map[string]bool)
	for _, r := range report.Results {
		for _, v := range r.Vulnerabilities {
			// empty
			if v.ID == "" {
				continue
			}

			// add only unique CVEs
			if _, ok := m[v.ID]; ok {
				continue
			}

			vul := &types.Vulnerability{
				ID:             v.ID,
				URL:            v.URL,
				Package:        v.Package,
				PackageVersion: v.PackageVersion,
				Title:          v.Title,
				Description:    v.Description,
				Severity:       v.Severity,
				Updated:        v.Updated,
				Image:          report.Image,
			}

			// filter
			if filter(vul) {
				continue
			}

			list = append(list, vul)
			m[v.ID] = true
		}
	}

	result := &types.VulnerabilityReport{
		Image:           report.Image,
		Vulnerabilities: list,
	}

	return result, nil
}
