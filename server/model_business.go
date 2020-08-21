package server

type Business struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Desc     string `json:"desc"`
	Qrcode   string `json:"qrcode"`
	Wxid     string `json:"wxid"`
	Status   int    `json:"status"`
	CreateAt string `json:"create_at"`
}

func (b *Business) Properties() []string {
	return []string{
		"id",
		"name",
		"avatar",
		"`desc`",
		"qrcode",
		"wxid",
		"status",
		"create_at",
	}
}

func (b *Business) Scan(s Scanner) error {
	return s(
		&b.Id,
		&b.Name,
		&b.Avatar,
		&b.Desc,
		&b.Qrcode,
		&b.Wxid,
		&b.Status,
		&b.CreateAt,
	)
}
