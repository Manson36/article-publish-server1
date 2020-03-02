package models

type PaginationCondition struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
