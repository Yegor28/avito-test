package entity

type Page struct {
	Page_number int `json:"page_number"`
	Page_size int   `json:"page_size"`
	Adverts []Ad    `json:"adverts"`
}
