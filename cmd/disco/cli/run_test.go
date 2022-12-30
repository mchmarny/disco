package cli

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestImgCmd(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping integration test")
	}

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runImagesCmd(c)
	assert.NoError(t, err)
}

func TestVulCmd(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping integration test")
	}

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runVulnsCmd(c)
	assert.NoError(t, err)
}

func TestLicCmd(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping integration test")
	}

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(newTestApp(t), set, nil)
	err := runLicenseCmd(c)
	assert.NoError(t, err)
}

func newTestApp(t *testing.T) *cli.App {
	app, err := newApp("v0.0.0-test", "test", time.Now().UTC().Format(time.RFC3339))
	assert.NoError(t, err)
	return app
}
