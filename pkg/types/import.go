package types

import "strings"

const (
	dataSetID              = "disco"
	tableLicenseID         = "licenses"
	tableVulnerabilityID   = "vulnerabilities"
	tablePackageID         = "packages"
	datasetLocationDefault = "US"

	TableKindUndefined TableKind = iota
	TableKindLicense
	TableKindVulnerability
	TableKindPackage

	TableKindUndefinedName     = "undefined"
	TableKindLicenseName       = "license"
	TableKindVulnerabilityName = "vulnerability"
	TableKindPackageName       = "package"

	SBOMFormatUndefined SBOMFormat = iota
	SBOMFormatSPDX
	SBOMFormatCycloneDX

	sbomFormatSPDXName      = "spdx"
	sbomFormatCycloneDXName = "cyclonedx"
)

type SBOMFormat int64

func (f SBOMFormat) String() string {
	switch f {
	case SBOMFormatSPDX:
		return sbomFormatSPDXName
	case SBOMFormatCycloneDX:
		return sbomFormatCycloneDXName
	default:
		return "undefined"
	}
}

func ParseSBOMFormat(s string) SBOMFormat {
	switch s {
	case sbomFormatSPDXName:
		return SBOMFormatSPDX
	case sbomFormatCycloneDXName:
		return SBOMFormatCycloneDX
	default:
		return SBOMFormatUndefined
	}
}

type TableKind int64

func (t TableKind) String() string {
	switch t {
	case TableKindLicense:
		return TableKindLicenseName
	case TableKindVulnerability:
		return TableKindVulnerabilityName
	case TableKindPackage:
		return TableKindPackageName
	default:
		return TableKindUndefinedName
	}
}

func NewPackageImportRequest(projectID, path, format string) *ImportRequest {
	return &ImportRequest{
		ProjectID:  projectID,
		TableKind:  TableKindPackage,
		Location:   datasetLocationDefault,
		DatasetID:  dataSetID,
		TableID:    tablePackageID,
		FilePath:   path,
		SBOMFormat: ParseSBOMFormat(format),
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
	ProjectID  string
	Location   string
	DatasetID  string
	TableKind  TableKind
	TableID    string
	FilePath   string
	SBOMFormat SBOMFormat
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
