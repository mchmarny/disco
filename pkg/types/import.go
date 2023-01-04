package types

import "strings"

const (
	dataSetID              = "disco"
	tableLicenseID         = "licenses"
	tableVulnerabilityID   = "vulnerabilities"
	datasetLocationDefault = "US"

	TableKindUndefined TableKind = iota
	TableKindLicense
	TableKindVulnerability

	TableKindUndefinedName     = "undefined"
	TableKindLicenseName       = "license"
	TableKindVulnerabilityName = "vulnerability"
)

type TableKind int64

func (t TableKind) String() string {
	switch t {
	case TableKindLicense:
		return TableKindLicenseName
	case TableKindVulnerability:
		return TableKindVulnerabilityName
	default:
		return TableKindUndefinedName
	}
}

func NewLicenseImportRequest(projectID, path string) *ImportRequest {
	return &ImportRequest{
		ProjectID: projectID,
		TableKind: TableKindLicense,
		Location:  datasetLocationDefault,
		DatasetID: dataSetID,
		TableID:   tableLicenseID,
		FilePath:  path,
	}
}

func NewVulnerabilityImportRequest(projectID, path string) *ImportRequest {
	return &ImportRequest{
		ProjectID: projectID,
		TableKind: TableKindVulnerability,
		Location:  datasetLocationDefault,
		DatasetID: dataSetID,
		TableID:   tableVulnerabilityID,
		FilePath:  path,
	}
}

type ImportRequest struct {
	ProjectID string
	Location  string
	DatasetID string
	TableKind TableKind
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
