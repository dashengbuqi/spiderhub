package helper

type ComboData struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type RequestParams struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Sort     string      `json:"sort"`
	Order    string      `json:"order"`
	Params   interface{} `json:"params"`
}
