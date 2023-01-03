package types

import "strings"

const (
	dataSetID            = "disco"
	tableLicenseID       = "licenses"
	tableVulnerabilityID = "vulnerabilities"
)

func NewLicenseImportRequest(projectID, path string) *ImportRequest {
	return &ImportRequest{
		ProjectID: projectID,
		DatasetID: dataSetID,
		TableID:   tableLicenseID,
		FilePath:  path,
	}
}

func NewVulnerabilityImportRequest(projectID, path string) *ImportRequest {
	return &ImportRequest{
		ProjectID: projectID,
		DatasetID: dataSetID,
		TableID:   tableVulnerabilityID,
		FilePath:  path,
	}
}

type ImportRequest struct {
	ProjectID string
	DatasetID string
	TableID   string
	FilePath  string
}

func ParseImageNameFromDigest(digest string) string {
	p := strings.Split(digest, "@")
	if len(p) > 0 {
		return p[0]
	}
	return digest
}

const imageDigestParts = 2

func ParseImageShaFromDigest(digest string) string {
	p := strings.Split(digest, "@")
	if len(p) == imageDigestParts {
		return p[1]
	}
	return ""
}
