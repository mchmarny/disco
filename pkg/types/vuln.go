package types

import "fmt"

type VulnerabilityReport struct {
	Image           string           `json:"image"`
	Vulnerabilities []*Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
	ID             string `json:"source"`
	Severity       string `json:"severity"`
	Package        string `json:"package"`
	PackageVersion string `json:"version"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	URL            string `json:"url"`
	Updated        string `json:"updated"`
}

func (v *Vulnerability) String() string {
	return fmt.Sprintf("ID: %s, URL: %s, Package: %s, Version: %s, Title: %s, Desc: %s, Severity: %s, Updated: %s", v.ID, v.URL, v.Package, v.PackageVersion, v.Title, v.Description, v.Severity, v.Updated)
}

func (v *Vulnerability) Hash() string {
	return string(Hash(v))
}
