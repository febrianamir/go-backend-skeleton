package response

type BasePaginateResponse struct {
	Total  uint `json:"-"`
	Offset uint `json:"-"`
	Page   uint `json:"-"`
	Limit  uint `json:"-"`
}
