package disco

import (
	"context"

	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/pkg/errors"
)

func DiscoverLicense(ctx context.Context, in *SimpleQuery) error {
	licenseFilter := func(v string) bool {
		return true
	}

	if err := scan(ctx, scanner.LicenseScanner, in, licenseFilter); err != nil {
		return errors.Wrap(err, "error scanning for vulnerabilities")
	}

	return nil
}
