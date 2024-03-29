package trivy

import (
	"fmt"
	"strings"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
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
			Provider:       parseCreatorInfo(report.CreationInfo),
			License:        parseLicense(p),
			Source:         parseSource(p),
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

const spdxToolKey = "Tool"

func parseCreatorInfo(in *spdx.CreationInfo) string {
	if in == nil {
		return ""
	}

	var sb strings.Builder

	for _, c := range in.Creators {
		if c.CreatorType == spdxToolKey {
			return c.Creator
		} else {
			sb.WriteString(c.Creator)
			sb.WriteString(" ")
		}
	}

	return strings.TrimSpace(sb.String())
}

func parseSource(in *spdx.Package) string {
	if in.PackageSourceInfo != "" {
		return in.PackageSourceInfo
	}

	for _, r := range in.PackageExternalReferences {
		if r.Locator != "" {
			return r.Locator
		}
	}

	return ""
}

func parseLicense(in *spdx.Package) string {
	if in.PackageLicenseDeclared != "" {
		return in.PackageLicenseDeclared
	}
	if in.PackageLicenseConcluded != "" {
		return in.PackageLicenseConcluded
	}
	if len(in.PackageLicenseInfoFromFiles) > 0 {
		return in.PackageLicenseInfoFromFiles[0]
	}

	return ""
}
