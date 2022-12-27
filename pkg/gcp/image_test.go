package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRequireImageFields(t *testing.T, img *ImageInfo) {
	assert.NotNil(t, img)
	assert.NotEmpty(t, img.Deployed, "Deployed")
	assert.NotEmpty(t, img.Name, "Name")
	assert.NotEmpty(t, img.Registry, "Registry")
	assert.NotEmpty(t, img.Project, "Project")
	assert.NotEmpty(t, img.ManifestURL(), "ManifestURL")
	assert.NotEmpty(t, img.URI(), "URI")
	assert.NotEmpty(t, img.URL(), "URL")
}

func TestImageInfoDeployed(t *testing.T) {
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/project/folder/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.Equal(t, "us-west1-docker.pkg.dev/project/folder/image", img.Deployed)
}

func TestARImageWithRegion(t *testing.T) {
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/project/folder/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsGCR)
	assert.True(t, img.IsAR)
	assert.Equal(t, "us-west1", img.Region)
	assert.Equal(t, "us-west1-docker.pkg.dev", img.Registry)
}

func TestARImageWithoutRegion(t *testing.T) {
	img, err := ParseImageInfo("us-docker.pkg.dev/project/folder/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsGCR)
	assert.True(t, img.IsAR)
	assert.Equal(t, "us", img.Region)
	assert.Equal(t, "us-docker.pkg.dev", img.Registry)
}

func TestGCRImageWithOutRegion(t *testing.T) {
	img, err := ParseImageInfo("gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.True(t, img.IsGCR)
	assert.False(t, img.IsAR)
	assert.Equal(t, "gcr.io", img.Registry)
}

func TestGCRImageWithRegion(t *testing.T) {
	img, err := ParseImageInfo("us.gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.True(t, img.IsGCR)
	assert.False(t, img.IsAR)
	assert.Equal(t, "us", img.Region)
	assert.Equal(t, "us.gcr.io", img.Registry)
}

func TestARImageProject(t *testing.T) {
	img, err := ParseImageInfo("us-docker.pkg.dev/project/folder/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsGCR)
	assert.True(t, img.IsAR)
	assert.Equal(t, "project", img.Project)
}

func TestGCRImageProject(t *testing.T) {
	img, err := ParseImageInfo("gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.True(t, img.IsGCR)
	assert.False(t, img.IsAR)
	assert.Equal(t, "project", img.Project)
}

func TestARImageFolderAndName(t *testing.T) {
	img, err := ParseImageInfo("us-docker.pkg.dev/project/folder/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsGCR)
	assert.True(t, img.IsAR)
	assert.Equal(t, "folder", img.Folder)
	assert.Equal(t, "image", img.Name)
}

func TestGCRImageFolderAndName(t *testing.T) {
	img, err := ParseImageInfo("gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.True(t, img.IsGCR)
	assert.False(t, img.IsAR)
	assert.Empty(t, img.Folder)
	assert.Equal(t, "image", img.Name)
}

func TestImageTag(t *testing.T) {
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/project/folder/name:v0.1.2")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsLatest)
	assert.Equal(t, "v0.1.2", img.Tag)
	img, err = ParseImageInfo("us-west1-docker.pkg.dev/project/folder/name")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsLatest)
	assert.Empty(t, img.Tag)
}

func TestImageLatest(t *testing.T) {
	img, err := ParseImageInfo("gcr.io/project/image:latest")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.True(t, img.IsLatest)
	assert.Empty(t, img.Tag)
	img, err = ParseImageInfo("gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsLatest)
	assert.Empty(t, img.Tag)
}

func TestImageDigest(t *testing.T) {
	img, err := ParseImageInfo("gcr.io/project/image@sha:123")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsLatest)
	assert.Empty(t, img.Tag)
	assert.Equal(t, "sha:123", img.Digest)
	img, err = ParseImageInfo("gcr.io/project/image")
	assert.NoError(t, err)
	testRequireImageFields(t, img)
	assert.False(t, img.IsLatest)
	assert.Empty(t, img.Tag)
	assert.Empty(t, img.Digest)
}

func TestImageManifestURL(t *testing.T) {
	_, err := ParseImageInfo("")
	assert.Error(t, err)
	img, err := ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "https://us-west1-docker.pkg.dev/v2/cloudy-demos/art/artomator/manifests/latest", img.ManifestURL())

	img, err = ParseImageInfo("us-west1-docker.pkg.dev/cloudy-demos/art/artomator:v0.8.3")
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.Equal(t, "https://us-west1-docker.pkg.dev/v2/cloudy-demos/art/artomator/manifests/v0.8.3", img.ManifestURL())
}
