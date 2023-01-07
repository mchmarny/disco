package bq

import (
	"time"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/rs/zerolog/log"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

func MakeSPDXPackageRows(in *v2_3.Document) []*PackageRow {
	list := make([]*PackageRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)
	batchID := time.Now().UTC().Unix()

	for _, p := range in.Packages {
		log.Info().Msgf("adding package %s from %s", p.PackageName, p.PackageSourceInfo)
		list = append(list, &PackageRow{
			BatchID:        batchID,
			Image:          types.ParseImageNameFromDigest(in.DocumentName),
			Sha:            types.ParseImageShaFromDigest(in.DocumentName),
			Format:         in.SPDXVersion,
			Provider:       types.SPDXCreatorInfo(in.CreationInfo),
			Originator:     p.PackageOriginator.Originator,
			Package:        p.PackageName,
			PackageVersion: p.PackageVersion,
			Source:         p.PackageSourceInfo,
			License:        p.PackageLicenseConcluded,
			Updated:        updated,
		})
	}

	return list
}
