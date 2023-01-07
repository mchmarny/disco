package scanner

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testImage = "alpine:latest"
)

var (
	emptyFilter = func(v interface{}) bool {
		return false
	}
)

func TestScan(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	fp := file.Name()

	cmd := exec.Command("touch", fp)
	err = runCmd(cmd, fp)
	assert.NoError(t, err)
}

func TestLicenses(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	dir, err := os.MkdirTemp("", "license")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	path := path.Join(dir, "test-license.json")

	rep, err := GetLicenses(testImage, path, emptyFilter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.NotEmpty(t, len(rep.Licenses))
}

func TestVulns(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	dir, err := os.MkdirTemp("", "vulns")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	path := path.Join(dir, "test-vuln.json")

	rep, err := GetVulnerabilities(testImage, path, emptyFilter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.NotEmpty(t, len(rep.Vulnerabilities))
}

func TestPackages(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	dir, err := os.MkdirTemp("", "packages")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	path := path.Join(dir, "test-pkg.json")

	rep, err := GetPackages(testImage, path, emptyFilter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.NotEmpty(t, len(rep.Packages))
}
