package links

type CreateLinkRequest struct {
	TargetURL   string  `json:"targetUrl"`
	CustomAlias *string `json:"customAlias"`
}

type ListLinksQuery struct {
	Page     int
	PageSize int
	Q        string
	Status   string
}

type ListLinksFilter struct {
	OwnerID        string
	Page           int
	PageSize       int
	Q              string
	Status         string
	IncludeDeleted bool
}

type LinkResponse struct {
	ID            string  `json:"id"`
	OwnerID       string  `json:"ownerId"`
	Code          string  `json:"code"`
	CustomAlias   *string `json:"customAlias,omitempty"`
	ShortURL      string  `json:"shortUrl"`
	TargetURL     string  `json:"targetUrl"`
	Status        string  `json:"status"`
	TotalClicks   int64   `json:"totalClicks"`
	LastClickedAt *string `json:"lastClickedAt,omitempty"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

type LinkListResponse struct {
	Items      []LinkResponse     `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}
