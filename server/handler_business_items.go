package server

import (
	"net/http"
	"strings"
)

func (s *Server) BusinessItemsGet(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()
	pager := s.GetPager(r)

	var bid int
	if err := s.DB.Query(
		"SELECT id FROM business_profile WHERE buid = ? AND status = 1 LIMIT 1",
		buid,
	).Scan(&bid); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	list := []*Item{}
	var item *Item
	if err := s.DB.Query(
		"SELECT "+strings.Join(item.Properties(), ",")+" "+
			"FROM item WHERE bid = ? "+pager.Sql(),
		bid,
	).ScanFunc(func(r Scanner) error {
		item = &Item{}
		if err := item.Scan(r); err != nil {
			return err
		}
		list = append(list, item)
		return nil
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, list)
}

func (s *Server) BusinessItemsPost(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()

	name := s.GetParam(r, "name").String()
	pic := s.GetParam(r, "pic").String()
	price := s.GetParam(r, "price").Int()
	desc := s.GetParam(r, "desc").String()

	if name == "" || price == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ITEM)
		return
	}

	var bid int
	if err := s.DB.Query(
		"SELECT id FROM business_profile WHERE buid = ? AND status = 1 LIMIT 1",
		buid,
	).Scan(&bid); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	if err := s.DB.Execute(
		"INSERT INTO item "+
			"(bid, name, price, pic, `desc`) "+
			"VALUES (?,?,?,?,?)",
		bid,
		name,
		price,
		pic,
		desc,
	).Error(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}

func (s *Server) BusinessItemsItidGet(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()
	itid := s.GetParam(r, "itid").Int()

	if itid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	var bid int
	if err := s.DB.Query(
		"SELECT id FROM business_profile WHERE buid = ? AND status = 1 LIMIT 1",
		buid,
	).Scan(&bid); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	item := &Item{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(item.Properties(), ", ")+" "+
			"FROM item WHERE bid = ? AND id = ? LIMIT 1",
		bid,
		itid,
	).ScanFunc(func(r Scanner) error {
		return item.Scan(r)
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, item)
}

func (s *Server) BusinessItemsItidPost(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()
	itid := s.GetParam(r, "itid").Int()

	if itid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	var bid int
	if err := s.DB.Query(
		"SELECT id FROM business_profile WHERE buid = ? AND status = 1 LIMIT 1",
		buid,
	).Scan(&bid); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	pic := s.GetParam(r, "pic").String()
	price := s.GetParam(r, "price").Int()
	desc := s.GetParam(r, "desc").String()
	status := s.GetParam(r, "status").Int()

	if price == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ITEM)
		return
	}

	if err := s.DB.Execute(
		"UPDATE item SET pic=?, price=?,`desc`=?,status=? WHERE id=? AND bid=? LIMIT 1",
		pic,
		price,
		desc,
		status,
		itid,
		bid,
	).Error(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}
