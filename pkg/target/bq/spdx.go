package bq

import (
	"time"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/rs/zerolog/log"
)

func MakePackageRows(in ...*types.PackageReport) []*PackageRow {
	list := make([]*PackageRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)
	batchID := time.Now().UTC().Unix()

	for _, r := range in {
		for _, p := range r.Packages {
			log.Debug().Msgf("adding package %s from %s", p.Package, p.Source)
			list = append(list, &PackageRow{
				BatchID:        batchID,
				Image:          types.ParseImageNameFromDigest(r.Image),
				Sha:            types.ParseImageShaFromDigest(r.Image),
				Format:         p.Format,
				Provider:       p.Provider,
				Package:        p.Package,
				PackageVersion: p.PackageVersion,
				Source:         p.Source,
				License:        p.License,
				Updated:        updated,
			})
		}
	}

	return list
}
