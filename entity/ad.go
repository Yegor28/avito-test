package entity

import "time"

type Ad struct {
	Id          int      `json: "ID"`
	Name        string   `json:"Name" validate:"required,lte=200"`
	Description string   `json:"Description" validate:"required,lte=1000"`
	Photos      []string `json:"Photos" validate:"required,lte=3"`
	Price       int      `json:"Price" validate:"required"`
	Time time.Time
}