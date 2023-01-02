package types

import (
	"errors"
)

// GetPager configures new Pager.
func GetPager[T any](list []T, pageSize int) (*Pager[T], error) {
	if list == nil {
		return nil, errors.New("empty list")
	}
	if pageSize < 1 {
		return nil, errors.New("page size must be a positive number")
	}
	return &Pager[T]{
		list:     list,
		pageSize: pageSize,
		page:     0,
	}, nil
}

// Pager pages through records.
type Pager[T any] struct {
	list     []T
	pageSize int
	page     int
}

// GetPageSize returns the configured page size.
func (p *Pager[T]) GetPageSize() int {
	return p.pageSize
}

// GetCurrentPage returns current page.
func (p *Pager[T]) GetCurrentPage() int {
	return p.page
}

// Reset resets the cursor back to it's initial stage.
func (p *Pager[T]) Reset() {
	p.page = 0
}

// Next returns next page from the list.
func (p *Pager[T]) Next() []T {
	start := p.page * p.pageSize
	stop := start + p.pageSize
	p.page++

	if p.page == 1 && p.pageSize >= len(p.list) {
		// one pager
		return p.list
	}

	if start >= len(p.list) {
		// reached end
		return nil
	}

	if stop > len(p.list) {
		// stop larger than the list, trim to size
		stop = len(p.list)
	}

	return p.list[start:stop]
}
