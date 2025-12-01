package httpx

import (
	"testing"
)

func TestGetPageAndLimit(t *testing.T){
	tests := []struct{
		page int
		limit int
		expectedOffset int32
		expectedLimit int32
		expectedCurrentPage int32
	}{
		{page: 1, limit: 10, expectedOffset: 0, expectedLimit: 10, expectedCurrentPage: 1},
		{page: 2, limit: 10, expectedOffset: 10, expectedLimit: 10, expectedCurrentPage: 2},
		{page: 0, limit: 10, expectedOffset: 0, expectedLimit: 10, expectedCurrentPage: 1},
		{page: -1, limit: 10, expectedOffset: 0, expectedLimit: 10, expectedCurrentPage: 1},
		{page: 1, limit: 0, expectedOffset: 0, expectedLimit: 10, expectedCurrentPage: 1},
		{page: 1, limit: -1, expectedOffset: 0, expectedLimit: 10, expectedCurrentPage: 1},
		{page: 1, limit: 101, expectedOffset: 0, expectedLimit: 100, expectedCurrentPage: 1},
		{page: 3, limit: 20, expectedOffset: 40, expectedLimit: 20, expectedCurrentPage: 3},
}
for _, tt := range tests {
	pagination := calculatePagination(int32(tt.page), int32(tt.limit))
	offset := pagination.Offset
	limit := pagination.Limit
	currentPage := pagination.CurrentPage

	if offset != tt.expectedOffset || limit != tt.expectedLimit || currentPage != tt.expectedCurrentPage {
		t.Errorf("For page %d and limit %d: expected offset %d, limit %d, current page %d but got offset %d, limit %d, current page %d",
			tt.page, tt.limit, tt.expectedOffset, tt.expectedLimit, tt.expectedCurrentPage,
			offset, limit, currentPage)
	}
}
}