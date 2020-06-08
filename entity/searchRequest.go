package entity

type SearchRequest struct {
	Limit      int    `json:"limit" default=10`
	OrderField string `json:"order_field"`
	// -1 по убыванию, 0 как встретилось, 1 по возрастанию
	OrderBy int `json:"order_by,omitempty"`
	Page int `json:"page"`
}