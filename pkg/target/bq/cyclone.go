package bq

import (
	"fmt"
	"strings"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/rs/zerolog/log"
)

func MakeCycloneDXPackageRows(in *cyclonedx.BOM) []*PackageRow {
	list := make([]*PackageRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)
	batchID := time.Now().UTC().Unix()

	for _, c := range *in.Components {
		log.Info().Msgf("adding package %s from %s", c.Name, c.BOMRef)
		list = append(list, &PackageRow{
			BatchID:        batchID,
			Image:          types.ParseImageNameFromDigest(in.Metadata.Component.Name),
			Sha:            types.ParseImageShaFromDigest(in.Metadata.Component.Name),
			Format:         in.BOMFormat,
			Provider:       cycloneDXCreatorInfo(in.Metadata),
			Package:        c.Name,
			PackageVersion: c.Version,
			Source:         c.BOMRef,
			License:        cycloneDXPackageLicense(c.Licenses),
			Updated:        updated,
		})
	}

	return list
}

func cycloneDXPackageLicense(in *cyclonedx.Licenses) string {
	if in == nil {
		return ""
	}

	for _, c := range *in {
		if c.License != nil {
			return c.License.ID
		}
		if c.Expression != "" {
			return c.Expression
		}
	}

	return ""
}

func cycloneDXCreatorInfo(in *cyclonedx.Metadata) string {
	if in == nil && in.Tools == nil {
		return ""
	}

	var sb strings.Builder

	for _, c := range *in.Tools {
		sb.WriteString(fmt.Sprintf("%s %s", c.Vendor, c.Version))
	}

	return strings.TrimSpace(sb.String())
}
