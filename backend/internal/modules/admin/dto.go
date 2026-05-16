package admin

import "tracklink/internal/modules/links"

type AdminLink links.LinkResponse

type PaginationResponse links.PaginationResponse

type AdminLinkListResponse struct {
	Items      []AdminLink        `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type AdminBlockLinkRequest struct {
	Reason *string `json:"reason"`
}

type ListLinksQuery struct {
	Page     int
	PageSize int
	Q        string
}

type ListLinksFilter struct {
	Page     int
	PageSize int
	Q        string
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}
