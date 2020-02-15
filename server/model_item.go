package server

type Item struct {
	Id       int    `json:"id"`
	Bid      int    `json:"bid"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Pic      string `json:"pic"`
	Desc     string `json:"desc"`
	Status   int    `json:"status"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

func (i *Item) Properties() []string {
	return []string{
		"id",
		"bid",
		"name",
		"price",
		"pic",
		"`desc`",
		"status",
		"create_at",
		"update_at",
	}
}

func (i *Item) Scan(s Scanner) error {
	return s(
		&i.Id,
		&i.Bid,
		&i.Name,
		&i.Price,
		&i.Pic,
		&i.Desc,
		&i.Status,
		&i.CreateAt,
		&i.UpdateAt,
	)
}
