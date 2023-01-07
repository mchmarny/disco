package disco

import (
	"context"
	"os/exec"
	"testing"

	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/source"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func setTestImplementations() {
	scanner.ScanVulnerability = func(digest, path string) *exec.Cmd {
		return exec.Command("cp", "../../etc/data/test-license.json", path) //nolint
	}

	scanner.ScanLicense = func(digest, path string) *exec.Cmd {
		return exec.Command("cp", "../../etc/data/test-vuln.json", path) //nolint
	}
}

func testImageProvider(ctx context.Context, in *types.ImagesQuery) ([]*types.ImageItem, error) {
	return []*types.ImageItem{
		{
			URI: "us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c",
			Context: map[string]interface{}{
				"project": "cloudrun",
				"folder":  "container",
				"name":    "hello",
			},
		},
	}, nil
}

func TestImage(t *testing.T) {
	source.ImageProvider = testImageProvider
	setTestImplementations()
	ctx := context.Background()

	err := DiscoverImages(ctx, nil)
	assert.Error(t, err, "error discovering images with nil query")

	err = DiscoverImages(ctx, &types.ImagesQuery{})
	assert.NoError(t, err, "error discovering images")

	err = DiscoverImages(ctx, &types.ImagesQuery{
		SimpleQuery: types.SimpleQuery{
			ProjectID: "test-project",
			OutputFmt: types.ParseOutputFormatOrDefault("yaml"),
		},
	})
	assert.NoError(t, err, "error discovering images with project")

	err = DiscoverImages(ctx, &types.ImagesQuery{
		URIOnly: true,
	})
	assert.NoError(t, err, "error discovering images with digests only")
}

func TestLicense(t *testing.T) {
	source.ImageProvider = testImageProvider
	setTestImplementations()
	ctx := context.Background()

	err := DiscoverLicenses(ctx, nil, nil)
	assert.Error(t, err, "error licenses images with nil query")

	err = DiscoverLicenses(ctx, &types.LicenseQuery{}, nil)
	assert.NoError(t, err, "error discovering license")

	err = DiscoverLicenses(ctx, &types.LicenseQuery{
		SimpleQuery: types.SimpleQuery{
			ProjectID:  "test-project",
			OutputPath: "../../tmp/license.tmp",
		},
	}, nil)
	assert.NoError(t, err, "error discovering license")
}

func TestVuln(t *testing.T) {
	source.ImageProvider = testImageProvider
	setTestImplementations()
	ctx := context.Background()

	err := DiscoverVulns(ctx, nil, nil)
	assert.Error(t, err, "error vulns images with nil query")

	err = DiscoverVulns(ctx, &types.VulnsQuery{}, nil)
	assert.NoError(t, err, "error discovering vulns")

	err = DiscoverVulns(ctx, &types.VulnsQuery{
		SimpleQuery: types.SimpleQuery{
			ProjectID: "test-project",
		},
	}, nil)
	assert.NoError(t, err, "error discovering vulns with project")
}

func TestFormatParse(t *testing.T) {
	f := types.ParseOutputFormatOrDefault("")
	assert.Equal(t, f, types.DefaultOutputFormat)
	f = types.ParseOutputFormatOrDefault("json")
	assert.Equal(t, f, types.JSONFormat)
	f = types.ParseOutputFormatOrDefault("yaml")
	assert.Equal(t, f, types.YAMLFormat)
}

func TestWriteOutput(t *testing.T) {
	err := writeOutput("", types.JSONFormat, nil)
	assert.Error(t, err, "error writing output with nil data")
	f := struct {
		Name string
	}{
		Name: "test",
	}
	err = writeOutput("", types.JSONFormat, f)
	assert.Nil(t, err, "error writing output with JSON format")
}
