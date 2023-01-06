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

	VulnReportFormatUndefined VulnReportFormat = iota
	VulnReportFormatDisco
	VulnReportFormatTrivy

	LicenseReportFormatUndefined LicenseReportFormat = iota
	LicenseReportFormatDisco
	LicenseReportFormatTrivy

	SBOMFormatSPDXName      = "spdx"
	SBOMFormatCycloneDXName = "cyclonedx"

	VulnReportFormatDiscoName = "disco"
	VulnReportFormatTrivyName = "trivy"

	LicenseReportFormatDiscoName = "disco"
	LicenseReportFormatTrivyName = "trivy"

	undefinedKind = "undefined"
)

type LicenseReportFormat int64

func (f LicenseReportFormat) String() string {
	switch f {
	case LicenseReportFormatDisco:
		return LicenseReportFormatDiscoName
	case LicenseReportFormatTrivy:
		return LicenseReportFormatTrivyName
	default:
		return undefinedKind
	}
}

func ParseLicenseReportFormat(s string) LicenseReportFormat {
	switch s {
	case LicenseReportFormatDiscoName:
		return LicenseReportFormatDisco
	case LicenseReportFormatTrivyName:
		return LicenseReportFormatTrivy
	default:
		return LicenseReportFormatUndefined
	}
}

type VulnReportFormat int64

func (f VulnReportFormat) String() string {
	switch f {
	case VulnReportFormatDisco:
		return VulnReportFormatDiscoName
	case VulnReportFormatTrivy:
		return VulnReportFormatTrivyName
	default:
		return undefinedKind
	}
}

func ParseVulnReportFormat(s string) VulnReportFormat {
	switch s {
	case VulnReportFormatDiscoName:
		return VulnReportFormatDisco
	case VulnReportFormatTrivyName:
		return VulnReportFormatTrivy
	default:
		return VulnReportFormatUndefined
	}
}

type SBOMFormat int64

func (f SBOMFormat) String() string {
	switch f {
	case SBOMFormatSPDX:
		return SBOMFormatSPDXName
	case SBOMFormatCycloneDX:
		return SBOMFormatCycloneDXName
	default:
		return undefinedKind
	}
}

func ParseSBOMFormat(s string) SBOMFormat {
	switch s {
	case SBOMFormatSPDXName:
		return SBOMFormatSPDX
	case SBOMFormatCycloneDXName:
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

func NewLicenseImportRequest(projectID, path, format string) *ImportRequest {
	return &ImportRequest{
		ProjectID:     projectID,
		TableKind:     TableKindLicense,
		Location:      datasetLocationDefault,
		DatasetID:     dataSetID,
		TableID:       tableLicenseID,
		FilePath:      path,
		LicenseFormat: ParseLicenseReportFormat(format),
	}
}

func NewVulnerabilityImportRequest(projectID, path, format string) *ImportRequest {
	return &ImportRequest{
		ProjectID:  projectID,
		TableKind:  TableKindVulnerability,
		Location:   datasetLocationDefault,
		DatasetID:  dataSetID,
		TableID:    tableVulnerabilityID,
		FilePath:   path,
		VulnFormat: ParseVulnReportFormat(format),
	}
}

type ImportRequest struct {
	ProjectID     string
	Location      string
	DatasetID     string
	TableKind     TableKind
	TableID       string
	FilePath      string
	SBOMFormat    SBOMFormat
	LicenseFormat LicenseReportFormat
	VulnFormat    VulnReportFormat
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
