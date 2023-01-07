package syft

import (
	"fmt"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	spdx "github.com/spdx/tools-golang/spdx/v2_3"
)

func ParsePackages(path string, filter types.ItemFilter) (*types.PackageReport, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	log.Debug().Msgf("parsing packages from %s", path)

	var report spdx.Document
	if err := types.UnmarshalFromFile(path, &report); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file: %s", path)
	}

	list := make([]*types.Package, 0)
	m := make(map[string]bool)
	for _, p := range report.Packages {
		// empty
		if p.PackageName == "" {
			continue
		}

		// add only unique licenses
		if _, ok := m[p.PackageName]; ok {
			continue
		}

		pkg := &types.Package{
			Package:        p.PackageName,
			PackageVersion: p.PackageVersion,
			Format:         report.SPDXVersion,
			Provider:       types.SPDXCreatorInfo(report.CreationInfo),
			Source:         p.PackageSourceInfo,
			License:        p.PackageLicenseConcluded,
		}

		if p.PackageOriginator != nil {
			pkg.Originator = p.PackageOriginator.Originator
		}

		// filter
		if filter(pkg) {
			continue
		}

		list = append(list, pkg)
		m[p.PackageName] = true
	}

	result := &types.PackageReport{
		Image:    report.DocumentName,
		Packages: list,
	}

	return result, nil
}
