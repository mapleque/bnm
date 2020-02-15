package server

type Address struct {
	Id       int    `json:"id"`
	Cuid     int    `json:"cuid"`
	Label    string `json:"label"`
	Reciever string `json:"reciever"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	UpdateAt string `json:"update_at"`
}

func (a *Address) Properties() []string {
	return []string{
		"id",
		"cuid",
		"label",
		"reciever",
		"address",
		"phone",
		"update_at",
	}
}

func (a *Address) Scan(s Scanner) error {
	return s(
		&a.Id,
		&a.Cuid,
		&a.Label,
		&a.Reciever,
		&a.Address,
		&a.Phone,
		&a.UpdateAt,
	)
}
