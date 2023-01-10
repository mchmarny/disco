package types

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	VulnSevUndefined VulnSev = iota
	VulnSevLow
	VulnSevMedium
	VulnSevHigh
	VulnSevCritical
)

type VulnSev int64 // Vulnerability severity

// String returns string representation of vulnerability severity.
func (s VulnSev) String() string {
	switch s {
	case VulnSevLow:
		return "low"
	case VulnSevMedium:
		return "medium"
	case VulnSevHigh:
		return "high"
	case VulnSevCritical:
		return "critical"
	default:
		return "undefined"
	}
}

// ParseMinVulnSeverityOrDefault parses vulnerability severity string.
func ParseMinVulnSeverityOrDefault(s string) VulnSev {
	switch strings.ToLower(s) {
	case "low":
		return VulnSevLow
	case "medium":
		return VulnSevMedium
	case "high":
		return VulnSevHigh
	case "critical":
		return VulnSevCritical
	default:
		return VulnSevUndefined
	}
}

type VulnsQuery struct {
	SimpleQuery
	CVE        string
	MinVulnSev VulnSev // Vulnerability severity
}

func (q *VulnsQuery) Validate() error {
	if err := q.SimpleQuery.Validate(); err != nil {
		return errors.Wrap(err, "invalid simple query")
	}

	if q.MinVulnSev != VulnSevUndefined && q.CVE != "" {
		return errors.New("min severity and CVE are mutually exclusive")
	}

	return nil
}

func (q *VulnsQuery) String() string {
	return fmt.Sprintf("project: %s, output: %s, format: %s, source: %s, uri: %s, cve: %s, severity: %s",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.ImageFile, q.ImageURI, q.CVE, q.MinVulnSev)
}

type VulnerabilityReport struct {
	Image           string            `json:"image"`
	Context         map[string]string `json:"context"`
	Vulnerabilities []*Vulnerability  `json:"vulnerabilities"`
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
	Image          string `json:"-"`
}

func (v *Vulnerability) String() string {
	return fmt.Sprintf("ID: %s, URL: %s, Package: %s, Version: %s, Title: %s, Desc: %s, Severity: %s, Updated: %s", v.ID, v.URL, v.Package, v.PackageVersion, v.Title, v.Description, v.Severity, v.Updated)
}
