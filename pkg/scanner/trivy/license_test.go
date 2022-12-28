package trivy

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestLicenseFilterParsing(t *testing.T) {
	expectedResults := 15

	filter := func(v interface{}) bool {
		return false
	}

	rep, err := testLicenseParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Licenses))
}

func testLicenseParsing(t *testing.T, filter types.ItemFilter) (*types.LicenseReport, error) {
	src := "../../../etc/data/test-license.json"
	img := "us-west1-docker.pkg.dev/cloudy-demos/artomator/artomator@sha256:b4a094e55244bc442bdaf2a5cd06a589f754ffc8ce09183868acaa79419cd88d"

	rep, err := ParseLicenses(img, src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
