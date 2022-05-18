package models

type Rectangle struct {
	N       float64 `json:"n"`
	S       float64 `json:"s"`
	E       float64 `json:"e"`
	W       float64 `json:"w"`
	Count   int64   `json:"count"`
	Opacity float64 `json:"opacity"`
}
