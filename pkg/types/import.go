package types

import (
	"strings"
)

const (
	TableKindUndefined TableKind = iota
	TableKindLicense
	TableKindVulnerability
	TableKindPackage

	TableKindUndefinedName     = "undefined"
	TableKindLicenseName       = "licenses"
	TableKindVulnerabilityName = "vulnerabilities"
	TableKindPackageName       = "packages"
)

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

type ImportRequest struct {
	ProjectID string
	Location  string
	DatasetID string
	TableKind TableKind
	TableID   string
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
