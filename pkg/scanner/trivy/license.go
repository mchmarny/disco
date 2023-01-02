package trivy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type licenseReport struct {
	Results []struct {
		Licenses []struct {
			PkgName string `json:"pkgName"`
			Name    string `json:"name"`
		} `json:"licenses"`
	} `json:"results"`
}

func ParseLicenses(image, path string, filter types.ItemFilter) (*types.LicenseReport, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	log.Debug().Msgf("parsing licenses from %s using filter %t", path, filter != nil)

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

			// add only unique licenses
			if _, ok := m[l.Name]; ok {
				continue
			}

			lic := &types.License{
				Name:   l.Name,
				Source: l.PkgName,
			}

			// filter
			if filter(lic) {
				continue
			}

			list = append(list, lic)
			m[l.Name] = true
		}
	}

	result := &types.LicenseReport{
		Image:    image,
		Licenses: list,
	}

	return result, nil
}
