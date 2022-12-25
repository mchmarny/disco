package trivy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
)

type licenseReport struct {
	Results []struct {
		Licenses []struct {
			PkgName string `json:"pkgName"`
			Name    string `json:"name"`
		} `json:"licenses"`
	} `json:"results"`
}

func ParseLicenses(image, path string, filters ...types.ItemFilter) (*types.LicenseReport, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", path)
	}

	var report licenseReport
	if err := json.Unmarshal(b, &report); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file: %s", path)
	}

	list := make([]*types.License, 0)
	m := make(map[string]bool)
	for _, r := range report.Results {
		for _, l := range r.Licenses {
			// empty
			if l.Name == "" {
				continue
			}
			// filter
			if filters == nil {
				for _, f := range filters {
					if f(l.Name) {
						continue
					}
				}
			}
			// add only unique licenses
			if _, ok := m[types.HashStr(l)]; ok {
				continue
			}
			list = append(list, &types.License{
				Name:   l.Name,
				Source: l.PkgName,
			})
		}
	}

	result := &types.LicenseReport{
		Image:    image,
		Licenses: list,
	}

	return result, nil
}
