package analysis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mchmarny/vctl/pkg/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	// occurrencesAPIBaseURL is the base URL for the Container Analysis API (projectID).
	occurrencesAPIBaseURL = "https://containeranalysis.googleapis.com/v1/projects/%s/occurrences?pageSize=500&filter="

	occurrencesAPIFilter = `resourceUrl="%s" AND kind="VULNERABILITY"`
)

type occurrenceList struct {
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

func GetImageVulnerabilities(ctx context.Context, projectID, imageURL string) ([]*Occurrence, error) {
	if imageURL == "" {
		return nil, errors.New("imageURL is empty")
	}
	if projectID == "" {
		return nil, errors.New("projectID is empty")
	}

	q := fmt.Sprintf(occurrencesAPIFilter, strings.TrimSpace(imageURL))
	log.Debug().Msgf("Query: '%s'.", q)

	u := fmt.Sprintf(occurrencesAPIBaseURL, projectID) + url.QueryEscape(q)
	log.Debug().Msgf("Encoded URL: '%s'.", u)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating image vulnerability request: %s", u)
	}

	var list occurrenceList
	if err := client.Request(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Occurrences, nil
}
