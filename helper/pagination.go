package helper

type Pagination struct {
	Total    int64 `json:"total"`
	PageSize int   `json:"page_size"`
	Page     int   `json:"page"`
}

func (this *Pagination) GetOffset() int {
	if this.Page == 0 {
		this.Page = 1
	}
	return (this.Page - 1) * this.PageSize
}

func (this *Pagination) GetLimit() int {
	return this.PageSize
}

func (this *Pagination) GetTotal() int64 {
	return this.Total
}
