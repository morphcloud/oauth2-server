package http_response

//easyjson:json
type SingleJSONResponse struct {
	Data interface{} `json:"data"`
}

//easyjson:json
type MultipleJSONResponse struct {
	Data interface{} `json:"data"`
}

//easyjson:json
type MultiplePaginatedJSONResponse struct {
	Data  interface{} `json:"data"`
	Links Links       `json:"links"`
	Meta  Meta        `json:"meta"`
}

//easyjson:json
type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

//easyjson:json
type Meta struct {
	CurrentPage uint64 `json:"currentPage"`
	From        uint64 `json:"from"`
	LastPage    uint64 `json:"lastPage"`
	Path        string `json:"path"`
	PerPage     uint64 `json:"perPage"`
	To          uint64 `json:"to"`
	Total       uint64 `json:"total"`
}
