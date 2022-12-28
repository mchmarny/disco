package types

const (
	KindImageName         = "image"
	KindLicenseName       = "license"
	KindVulnerabilityName = "vulnerability"
	KindUndefinedName     = "undefined"

	KindUndefined DiscoKind = iota
	KindImage
	KindLicense
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
