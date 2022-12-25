package trivy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
)

type vulnerabilityReport struct {
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

func ParseVulnerabilities(image, path string, filters ...types.ItemFilter) (*types.VulnerabilityReport, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", path)
	}

	var report vulnerabilityReport
	if err := json.Unmarshal(b, &report); err != nil {
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
			// filter
			if filters == nil {
				for _, f := range filters {
					if f(v.ID) {
						continue
					}
				}
			}
			// add only unique licenses
			if _, ok := m[types.ToHash(v.ID, v.Package, v.PackageVersion)]; ok {
				continue
			}
			list = append(list, &types.Vulnerability{
				ID:             v.ID,
				URL:            v.URL,
				Package:        v.Package,
				PackageVersion: v.PackageVersion,
				Title:          v.Title,
				Description:    v.Description,
				Severity:       v.Severity,
				Updated:        v.Updated,
			})
		}
	}

	result := &types.VulnerabilityReport{
		Image:           image,
		Vulnerabilities: list,
	}

	return result, nil
}
