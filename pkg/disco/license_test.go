package disco

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func testFilterOnList[T any](filter types.ItemFilter, list []*T) int {
	found := 0
	for _, v := range list {
		if !filter(v) {
			found++
		}
	}
	return found
}

func TestLicenseNoFilter(t *testing.T) {
	filter := makeLicenseFilter(&types.LicenseQuery{})

	list := []*types.License{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
		{
			Name: "test",
		},
	}

	found := testFilterOnList(filter, list)
	assert.Equal(t, 3, found)
}

func TestLicenseFilter(t *testing.T) {
	filter := makeLicenseFilter(&types.LicenseQuery{
		TypeFilter: "test",
	})

	list := []*types.License{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
		{
			Name: "test",
		},
	}

	found := testFilterOnList(filter, list)
	assert.Equal(t, 1, found)
}
