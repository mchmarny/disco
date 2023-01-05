package disco

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestVulnImport(t *testing.T) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "cloudy-demos", "test")
	set.String("dataset", "disco_test", "test")
	set.String("file", "../../../etc/data/test-vuln.json", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runVulnImportCmd(c)
	assert.NoError(t, err)
}

func TestLicenseImport(t *testing.T) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "cloudy-demos", "test")
	set.String("dataset", "disco_test", "test")
	set.String("file", "../../../etc/data/test-license.json", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runLicenseImportCmd(c)
	assert.NoError(t, err)
}

func TestPackageSPDXImport(t *testing.T) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "cloudy-demos", "test")
	set.String("dataset", "disco_test", "test")
	set.String("format", "spdx", "test")
	set.String("file", "../../../etc/data/spdx23.json", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runPackageImportCmd(c)
	assert.NoError(t, err)
}

func TestPackageCycloneDXImport(t *testing.T) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "cloudy-demos", "test")
	set.String("dataset", "disco_test", "test")
	set.String("format", "cyclonedx", "test")
	set.String("file", "../../../etc/data/cyclonedx12.json", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runPackageImportCmd(c)
	assert.NoError(t, err)
}
