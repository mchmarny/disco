package project

import (
	"context"
	"net/http"

	"github.com/mchmarny/vctl/pkg/client"
	"github.com/pkg/errors"
)

const (
	projectAPIBaseURL = "https://cloudresourcemanager.googleapis.com/v1/projects"

	// ProjectStateActive is the state of active project.
	ProjectStateActive = "ACTIVE"
)

type projectList struct {
	Projects []*Project `json:"projects"`
}

// Project represents GCP project.
type Project struct {
	ProjectNumber  string `json:"projectNumber"`
	ProjectID      string `json:"projectId"`
	LifecycleState string `json:"lifecycleState"`
}

// GetProjects returns list of projects available to the current user.
func GetProjects(ctx context.Context) ([]*Project, error) {
	req, err := http.NewRequest(http.MethodGet, projectAPIBaseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating request")
	}

	var list projectList
	if err := client.Request(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Projects, nil
}
