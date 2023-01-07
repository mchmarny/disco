package types

const (
	KindImageName         = "image"
	KindLicenseName       = "license"
	KindVulnerabilityName = "vulnerability"
	KindPackageName       = "package"
	KindUndefinedName     = "undefined"

	KindUndefined DiscoKind = iota
	KindImage
	KindLicense
	KindPackage
	KindVulnerability
)

type DiscoKind int64

func (k DiscoKind) String() string {
	switch k {
	case KindImage:
		return KindImageName
	case KindLicense:
		return KindLicenseName
	case KindVulnerability:
		return KindVulnerabilityName
	case KindPackage:
		return KindPackageName
	default:
		return KindUndefinedName
	}
}

type Meta struct {
	Kind    string `json:"kind"`
	Version string `json:"version"`
	Created string `json:"created"`
	Count   *int   `json:"count,omitempty"`
}
