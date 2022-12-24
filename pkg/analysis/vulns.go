package analysis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mchmarny/vctl/pkg/client"
	"github.com/pkg/errors"
)

const (
	// occurrencesAPIBaseURL is the base URL for the Container Analysis API (projectID).
	occurrencesAPIBaseURL = "https://containeranalysis.googleapis.com/v1/projects/%s/occurrences?pageSize=500&filter="

	occurrencesAPIFilter = `resourceUrl="%q" kind="VULNERABILITY"`
)

type occurrence struct {
	Occurrences []*Occurrence `json:"occurrences,omitempty"`
}

// Occurrence represents a single vulnerability occurrence.
type Occurrence []struct {
	Name          string `json:"name,omitempty"`
	URI           string `json:"resourceUri,omitempty"`
	Note          string `json:"noteName,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Updated       string `json:"updateTime,omitempty"`
	Vulnerability struct {
		Severity    string `json:"severity,omitempty"`
		Score       int    `json:"cvssScore,omitempty"`
		Description string `json:"shortDescription,omitempty"`
		RelatedURLs []struct {
			URL   string `json:"url,omitempty"`
			Label string `json:"label,omitempty"`
		} `json:"relatedUrls,omitempty"`
		Packages []struct {
			Name              string `json:"affectedPackage"`
			Affected          string `json:"affectedCpeUri"`
			Type              string `json:"packageType"`
			EffectiveSeverity string `json:"effectiveSeverity"`
			Version           struct {
				Name     string `json:"name,omitempty"`
				Revision string `json:"revision,omitempty"`
				Kind     string `json:"kind,omitempty"`
				FullName string `json:"fullName,omitempty"`
			} `json:"affectedVersion,omitempty"`
		} `json:"packageIssue,omitempty"`
		CVSSv3 struct {
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

func GetImageVulnerabilities(ctx context.Context, projectID, digest string) ([]*Occurrence, error) {
	if digest == "" {
		return nil, errors.New("project number is empty")
	}
	if projectID == "" {
		return nil, errors.New("projectID is empty")
	}

	q := fmt.Sprintf(occurrencesAPIFilter, digest)
	u := fmt.Sprintf(occurrencesAPIBaseURL, projectID) + url.QueryEscape(q)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating image vulnerability request: %s", u)
	}

	var list occurrence
	if err := client.Request(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Occurrences, nil
}
