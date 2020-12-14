package responses

type DataRes struct {
	APIVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type Responser interface {
	Response() (res interface{})
}

type BinaryResponser interface {
	GetBinary() (res []byte)
}

type Paging struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
