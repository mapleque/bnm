package server

import (
	"encoding/json"
)

type Order struct {
	Id         int    `json:"id"`
	Bid        int    `json:"bid"`
	Cuid       int    `json:"cuid"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Counts     int    `json:"counts"`
	Reciever   string `json:"reciever"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Stage      string `json:"stage"`
	Status     int    `json:"status"`
	Additional string `json:"additional"`
	ExpNo      string `json:"exp_no"`
	UpdateAt   string `json:"update_at"`
	CreateAt   string `json:"create_at"`
}

func (o *Order) Properties() []string {
	return []string{
		"id",
		"bid",
		"cuid",
		"name",
		"price",
		"counts",
		"reciever",
		"address",
		"phone",
		"stage",
		"status",
		"additional",
		"exp_no",
		"update_at",
		"create_at",
	}
}

func (o *Order) Scan(s Scanner) error {
	return s(
		&o.Id,
		&o.Bid,
		&o.Cuid,
		&o.Name,
		&o.Price,
		&o.Counts,
		&o.Reciever,
		&o.Address,
		&o.Phone,
		&o.Stage,
		&o.Status,
		&o.Additional,
		&o.ExpNo,
		&o.UpdateAt,
		&o.CreateAt,
	)
}

func (o *Order) String() string {
	bt, _ := json.Marshal(o)
	return string(bt)
}
