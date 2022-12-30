package run

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	projectAPIBaseURL = "https://cloudresourcemanager.googleapis.com/v1/projects"

	// ProjectStateActive is the state of active project.
	ProjectStateActive = "ACTIVE"
)

type projectList struct {
	Projects []*project `json:"projects"`
}

// Project represents GCP project.
type project struct {
	Number string `json:"projectNumber"`
	ID     string `json:"projectId"`
	State  string `json:"lifecycleState"`
}

func getProjects(ctx context.Context) ([]*project, error) {
	req, err := http.NewRequest(http.MethodGet, projectAPIBaseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating request")
	}

	var list projectList
	if err := apiClient.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Projects, nil
}

func isQualifiedProject(ctx context.Context, p *project, filterID string) bool {
	log.Debug().Msgf("qualifying project: %s (%s - %s)", p.ID, p.Number, p.State)

	if filterID != "" && p.ID == filterID {
		log.Debug().Msgf("skipping: %s (filter: %s)", p.ID, filterID)
		return false
	}

	if p.State != ProjectStateActive {
		log.Debug().Msgf("skipping: %s (inactive)", p.ID)
		return false
	}

	on, err := isAPIEnabled(ctx, p.Number, cloudRunAPI)
	if err != nil {
		log.Error().Err(err).Msgf("error checking Cloud Run API: %s", p.ID)
		return false
	}

	if !on {
		log.Debug().Msgf("skipping: %s (API not enabled)", p.ID)
		return false
	}

	return true
}
