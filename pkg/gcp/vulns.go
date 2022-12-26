package gcp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	// occurrencesAPIBaseURL is the base URL for the Container Analysis API (projectID).
	occurrencesAPIBaseURL = "https://containeranalysis.googleapis.com/v1/projects/%s/occurrences?pageSize=500&filter="

	occurrencesAPIFilter = `resourceUrl="%s" AND kind="VULNERABILITY"`

	// https://cloud.google.com/container-analysis/docs/go-scanning-automatically
	cveAPIFilter = `noteId="%s"`

	/*
		kind="PACKAGE"
		noteProjectId="goog-analysis"
		resourceUrl="projects/PROJECT_ID/locations/LOCATION/repositories/REPOSITORY/mavenArtifacts/PACKAGE_NAME:VERSION"
	*/

)

// Vulnerability represents a single vulnerability.
type OccurrenceList struct {
	Occurrences []*Occurrence `json:"occurrences,omitempty"`
}

// Occurrence represents a single vulnerability occurrence.
type Occurrence struct {
	Name          string `json:"name,omitempty"`
	URI           string `json:"resourceUri,omitempty"`
	Updated       string `json:"updateTime,omitempty"`
	Vulnerability struct {
		Severity     string  `json:"severity,omitempty"`
		CVSSScore    float64 `json:"cvssScore,omitempty"`
		PackageIssue []struct {
			AffectedPackage string `json:"affectedPackage,omitempty"`
			AffectedVersion struct {
				Name     string `json:"name,omitempty"`
				Revision string `json:"revision,omitempty"`
				Kind     string `json:"kind,omitempty"`
				FullName string `json:"fullName,omitempty"`
			} `json:"affectedVersion,omitempty"`
			Fixed        string `json:"fixedPackage,omitempty"`
			FixedVersion struct {
				Name     string `json:"name,omitempty"`
				Revision string `json:"revision,omitempty"`
				Kind     string `json:"kind,omitempty"`
				FullName string `json:"fullName,omitempty"`
			} `json:"fixedVersion,omitempty"`
			Type              string `json:"packageType"`
			EffectiveSeverity string `json:"effectiveSeverity"`
		} `json:"packageIssue,omitempty"`
		ShortDescription string `json:"shortDescription,omitempty"`
		RelatedURLs      []struct {
			URL string `json:"url,omitempty"`
		} `json:"relatedUrls,omitempty"`
		EffectiveSeverity string `json:"effectiveSeverity"`
		CVSSv3            struct {
			BaseScore             float64 `json:"baseScore,omitempty"`
			EmployabilityScore    float64 `json:"exploitabilityScore,omitempty"`
			ImpactScore           float64 `json:"impactScore,omitempty"`
			AttackVector          string  `json:"attackVector,omitempty"`
			AttackComplexity      string  `json:"attackComplexity,omitempty"`
			PrivilegesRequired    string  `json:"privilegesRequired,omitempty"`
			UserInteraction       string  `json:"userInteraction,omitempty"`
			Scope                 string  `json:"scope,omitempty"`
			ConfidentialityImpact string  `json:"confidentialityImpact,omitempty"`
			IntegrityImpact       string  `json:"integrityImpact,omitempty"`
			AvailabilityImpact    string  `json:"availabilityImpact,omitempty"`
		} `json:"cvssv3,omitempty"`
	} `json:"vulnerability,omitempty"`
}

// GetCVEVulnerabilities returns all vulnerabilities for a given CVE.
func GetCVEVulnerabilities(ctx context.Context, projectID, cveID string) ([]*Occurrence, error) {
	return getVulnerabilities(ctx, projectID, "", cveID)
}

// GetImageVulnerabilities returns all vulnerabilities for a given image.
func GetImageVulnerabilities(ctx context.Context, projectID, imageURL string) ([]*Occurrence, error) {
	return getVulnerabilities(ctx, projectID, imageURL, "")
}

func getVulnerabilities(ctx context.Context, projectID, imageURL, cveID string) ([]*Occurrence, error) {
	if projectID == "" {
		return nil, errors.New("projectID is empty")
	}
	if imageURL == "" && cveID == "" {
		return nil, errors.New("either CVE or image are required")
	}

	u := fmt.Sprintf(occurrencesAPIBaseURL, projectID)
	q := ""
	if cveID != "" {
		q = fmt.Sprintf(cveAPIFilter, strings.TrimSpace(cveID))
	} else {
		q = fmt.Sprintf(occurrencesAPIFilter, strings.TrimSpace(imageURL))
	}
	log.Debug().Msgf("query: '%s'", q)
	u += url.QueryEscape(q)
	log.Debug().Msgf("encoded URL: '%s'", u)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating image vulnerability request: %s", u)
	}

	var list OccurrenceList
	if err := APIClient.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	if len(list.Occurrences) == 0 {
		return make([]*Occurrence, 0), nil
	}

	return list.Occurrences, nil
}
