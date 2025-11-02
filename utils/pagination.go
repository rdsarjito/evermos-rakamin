package utils

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"limit" query:"limit"`
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationParams extracts pagination params from request with defaults
func GetPaginationParams(page, limit int) (offset int, limitVal int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // max limit
	}

	offset = (page - 1) * limit
	return offset, limit
}

// GetTotalPages calculates total pages from total records and limit
func GetTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	return totalPages
}

