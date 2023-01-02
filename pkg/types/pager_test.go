package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagerWithInt64(t *testing.T) {
	list := []int64{0, 1, 2, 3, 4, 5, 6, 7}

	t.Run("nil list", func(t *testing.T) {
		_, err := GetPager[int64](nil, 10)
		assert.Error(t, err)
	})

	t.Run("page size", func(t *testing.T) {
		_, err := GetPager(list, 0)
		assert.Error(t, err)
	})

	t.Run("page larger than list", func(t *testing.T) {
		p, err := GetPager(list, 100)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		n := p.Next()
		assert.Equal(t, len(list), len(n))
	})

	t.Run("page equal to list", func(t *testing.T) {
		p, err := GetPager(list, len(list))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		n := p.Next()
		assert.Equal(t, len(list), len(n))
	})

	t.Run("page loop check", func(t *testing.T) {
		p, err := GetPager(list, len(list))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, 0, p.GetCurrentPage())

		n := p.Next()
		assert.Equal(t, 1, p.GetCurrentPage())
		assert.Equal(t, len(list), len(n))

		n = p.Next()
		assert.Equal(t, 2, p.GetCurrentPage())
		assert.Nil(t, n)

		n = p.Next()
		assert.Equal(t, 3, p.GetCurrentPage())
		assert.Nil(t, n)
	})

	t.Run("page end to end", func(t *testing.T) {
		p, err := GetPager(list, 3)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		s := p.GetPageSize()
		assert.Equal(t, 3, s)
		n := p.Next()
		assert.Equal(t, 3, len(n))
		assert.Equal(t, int64(0), n[0])
		assert.Equal(t, int64(1), n[1])
		assert.Equal(t, int64(2), n[2])
		n = p.Next()
		assert.Equal(t, 3, len(n))
		assert.Equal(t, int64(3), n[0])
		assert.Equal(t, int64(4), n[1])
		assert.Equal(t, int64(5), n[2])
		n = p.Next() // gets last 2
		assert.Equal(t, 2, len(n))
		assert.Equal(t, int64(6), n[0])
		assert.Equal(t, int64(7), n[1])
		n = p.Next() // no more left
		assert.Nil(t, n)
		p.Reset()
		n = p.Next()
		assert.Equal(t, 3, len(n))
		assert.Equal(t, int64(0), n[0])
		assert.Equal(t, int64(1), n[1])
		assert.Equal(t, int64(2), n[2])
	})
}

type PageTestStruct struct {
	ID  int64
	Val interface{}
}

func TestPagerWithStruct(t *testing.T) {
	list := make([]PageTestStruct, 0)
	for i := 0; i < 10; i++ {
		list = append(list, PageTestStruct{
			ID:  int64(i),
			Val: i,
		})
	}

	t.Run("struct page end to end", func(t *testing.T) {
		p, err := GetPager(list, 4)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		s := p.GetPageSize()
		assert.Equal(t, 4, s)
		n := p.Next()
		assert.Equal(t, 4, len(n))
		assert.Equal(t, int64(0), n[0].ID)
		assert.Equal(t, int64(1), n[1].ID)
		assert.Equal(t, int64(2), n[2].ID)
		assert.Equal(t, int64(3), n[3].ID)
		n = p.Next()
		assert.Equal(t, 4, len(n))
		assert.Equal(t, int64(4), n[0].ID)
		assert.Equal(t, int64(5), n[1].ID)
		assert.Equal(t, int64(6), n[2].ID)
		assert.Equal(t, int64(7), n[3].ID)
		n = p.Next() // gets last 2
		assert.Equal(t, 2, len(n))
		assert.Equal(t, int64(8), n[0].ID)
		assert.Equal(t, int64(9), n[1].ID)
		n = p.Next() // no more left
		assert.Nil(t, n)
		p.Reset()
		n = p.Next()
		assert.Equal(t, 4, len(n))
		assert.Equal(t, int64(0), n[0].ID)
		assert.Equal(t, int64(1), n[1].ID)
		assert.Equal(t, int64(2), n[2].ID)
		assert.Equal(t, int64(3), n[3].ID)
	})
}
