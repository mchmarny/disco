package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	client = &testAPIClient{}
	expectedProjects := 4
	list, err := getProjects(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedProjects, len(list))
}

func TestProjectFilter(t *testing.T) {
	client = &testAPIClient{}

	assert.False(t, isQualifiedProject(context.Background(), &project{
		ID:     "test",
		State:  ProjectStateActive,
		Number: "799736955886",
	}, "test"))

	assert.False(t, isQualifiedProject(context.Background(), &project{
		ID:     "test",
		State:  "SUSPENDED",
		Number: "799736955886",
	}, "test"))

	assert.True(t, isQualifiedProject(context.Background(), &project{
		ID:     "test",
		State:  ProjectStateActive,
		Number: "799736955886",
	}, "test2"))
}
