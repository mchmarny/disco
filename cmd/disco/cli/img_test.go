package cli

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestImgCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runImagesCmd(c)
	assert.NoError(t, err)
}

func TestVulCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String(
		"image",
		"gcr.io/cloudy-demos/system-package@sha256:92464e03ea61e922c46040359b32f55b7f38bb20eaa25b9679338fa07e5a71d7",
		"test",
	)

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runVulnsCmd(c)
	assert.NoError(t, err)
}

func TestLicCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String(
		"file",
		"../../../data/images.txt",
		"test",
	)

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runLicenseCmd(c)
	assert.NoError(t, err)
}

func TestPkgCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String(
		"file",
		"../../../data/images.txt",
		"test",
	)

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runPackageCmd(c)
	assert.NoError(t, err)
}

func newTestApp(t *testing.T) *cli.App {
	app, err := newApp("v0.0.0-test", "test", time.Now().UTC().Format(time.RFC3339))
	assert.NoError(t, err)
	return app
}
