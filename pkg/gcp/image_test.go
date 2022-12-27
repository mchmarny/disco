package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: refactor to declarative test cases.

func TestImageManifestURL(t *testing.T) {
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "https://us-west1-docker.pkg.dev/v2/cloudy-demos/art/artomator/manifests/latest", img.ManifestURL())

	img, err = ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator:v0.8.3")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "https://us-west1-docker.pkg.dev/v2/cloudy-demos/art/artomator/manifests/v0.8.3", img.ManifestURL())
}

func TestParseARImage(t *testing.T) {
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.False(t, img.IsGCR)
	assert.True(t, img.IsAR)
	assert.False(t, img.IsLatest)
	assert.Equal(t, "us-west1-docker.pkg.dev/cloudy-demos/art/artomator", img.Deployed)
	assert.Equal(t, "us-west1", img.Region)
	assert.Equal(t, "us-west1-docker.pkg.dev", img.Registry)
	assert.Equal(t, "cloudy-demos", img.Project)
	assert.Equal(t, "art", img.Folder)
	assert.Equal(t, "artomator", img.Name)
	assert.Equal(t, "latest", img.Tag)
	assert.Equal(t, "", img.Digest)

	img, err = ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator:latest")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.IsLatest)
	assert.Equal(t, "art", img.Folder)
	assert.Equal(t, "artomator", img.Name)
	assert.Equal(t, "latest", img.Tag)
	assert.Equal(t, "", img.Digest)

	img, err = ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator:v0.8.3")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.False(t, img.IsLatest)
	assert.Equal(t, "us-west1", img.Region)
	assert.Equal(t, "art", img.Folder)
	assert.Equal(t, "artomator", img.Name)
	assert.Equal(t, "v0.8.3", img.Tag)
	assert.Equal(t, "", img.Digest)

	img, err = ParseImageInfo("us-docker.pkg.dev/cloudy-demos/art/artomator@sha256:1234567890")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.False(t, img.IsLatest)
	assert.Equal(t, "us", img.Region)
	assert.Equal(t, "art", img.Folder)
	assert.Equal(t, "artomator", img.Name)
	assert.Equal(t, "", img.Tag)
	assert.Equal(t, "sha256:1234567890", img.Digest)
}

func TestParseGCRImage(t *testing.T) {
	_, err := ParseImageInfo("invalid.url")
	assert.Error(t, err)

	img, err := ParseImageInfo("gcr.io/cloudy-demos/hello-broken")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.IsGCR)
	assert.False(t, img.IsAR)
	assert.False(t, img.IsLatest)
	assert.Equal(t, "gcr.io/cloudy-demos/hello-broken", img.Deployed)
	assert.Equal(t, "", img.Region)
	assert.Equal(t, "gcr.io", img.Registry)
	assert.Equal(t, "cloudy-demos", img.Project)
	assert.Equal(t, "", img.Folder)
	assert.Equal(t, "hello-broken", img.Name)
	assert.Equal(t, "latest", img.Tag)
	assert.Empty(t, img.Digest)

	img, err = ParseImageInfo("gcr.io/cloudy-demos/hello-broken:latest")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.IsLatest)
	assert.Equal(t, "", img.Region)
	assert.Equal(t, "hello-broken", img.Name)
	assert.Equal(t, "latest", img.Tag)

	img, err = ParseImageInfo("gcr.io/cloudy-demos/hello-broken@sha256:1234567890")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "", img.Region)
	assert.Equal(t, "hello-broken", img.Name)
	assert.Equal(t, "", img.Tag)
	assert.Equal(t, "sha256:1234567890", img.Digest)

	img, err = ParseImageInfo("us.gcr.io/cloudy-demos/hello-broken:v0.8.3")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "us", img.Region)
	assert.Equal(t, "us.gcr.io", img.Registry)
	assert.Equal(t, "hello-broken", img.Name)
	assert.Equal(t, "v0.8.3", img.Tag)
	assert.Empty(t, img.Digest)
	assert.NotEmpty(t, img.URL)
	assert.NotEmpty(t, img.URI)
}

func TestParseName(t *testing.T) {
	info := &ImageInfo{}
	k := parseName("artomator", info)
	assert.True(t, k)
	assert.Equal(t, "artomator", info.Name)
	assert.Equal(t, "latest", info.Tag)
	assert.Empty(t, info.Digest)
	assert.False(t, info.IsGCR)
	assert.False(t, info.IsAR)
	assert.False(t, info.IsLatest)

	info = &ImageInfo{}
	k = parseName("artomator:latest", info)
	assert.True(t, k)
	assert.Equal(t, "artomator", info.Name)
	assert.Equal(t, "latest", info.Tag)
	assert.Empty(t, info.Digest)
	assert.False(t, info.IsGCR)
	assert.False(t, info.IsAR)
	assert.True(t, info.IsLatest)

	info = &ImageInfo{}
	k = parseName("artomator:v0.8.3", info)
	assert.True(t, k)
	assert.Equal(t, "artomator", info.Name)
	assert.Equal(t, "v0.8.3", info.Tag)
	assert.Empty(t, info.Digest)
	assert.False(t, info.IsGCR)
	assert.False(t, info.IsAR)

	info = &ImageInfo{}
	k = parseName("artomator@sha256:1234567890", info)
	assert.True(t, k)
	assert.Equal(t, "artomator", info.Name)
	assert.Equal(t, "sha256:1234567890", info.Digest)
	assert.Empty(t, info.Tag)
	assert.False(t, info.IsGCR)
	assert.False(t, info.IsAR)
}
