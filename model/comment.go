package model

type Comment struct {
	ID      int  `json:"id"`
	Bought  bool `json:"bought"`
	ImgPath string
}
