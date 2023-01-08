package disco

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPackageFilter(t *testing.T) {
	filter := makePackageFilter(&types.PackageQuery{
		NamePart: "ES",
	})

	list := []*types.Package{
		{
			Package: "test",
		},
		{
			Package: "go",
		},
		{
			Package: "foo",
		},
	}

	found := testFilterOnList(filter, list)
	assert.Equal(t, 1, found)
}
