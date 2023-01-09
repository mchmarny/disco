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
	return parseImageDigestParts(digest, 0)
}

func ParseImageShaFromDigest(digest string) string {
	return parseImageDigestParts(digest, 1)
}

func ParseImageShaFromDigestWithoutPrefix(digest string) string {
	s := parseImageDigestParts(digest, 1)
	return strings.TrimPrefix(s, "sha256:")
}

func parseImageDigestParts(digest string, index int) string {
	if index < 0 {
		return ""
	}

	p := strings.Split(digest, "@")
	if index < len(p) {
		return p[index]
	}

	return ""
}
