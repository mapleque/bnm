package server

import (
	"net/http"
	"strings"
)

func (s *Server) CustomerBusinessBidItemsGet(w http.ResponseWriter, r *http.Request) {
	bid := s.GetParam(r, "bid").Int()
	if bid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	pager := s.GetPager(r)

	list := []*Item{}
	var item *Item
	if err := s.DB.Query(
		"SELECT "+strings.Join(item.Properties(), ",")+" "+
			"FROM item WHERE bid = ? AND status = 1 "+pager.Sql(),
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

func (s *Server) CustomerBusinessBidItemsItidGet(w http.ResponseWriter, r *http.Request) {
	bid := s.GetParam(r, "bid").Int()
	itid := s.GetParam(r, "itid").Int()
	if bid == 0 || itid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	item := &Item{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(item.Properties(), ", ")+" "+
			"FROM item WHERE bid = ? AND status = 1 AND id = ? LIMIT 1",
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

func (s *Server) CustomerBusinessBidItemsItidPost(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	bid := s.GetParam(r, "bid").Int()
	itid := s.GetParam(r, "itid").Int()
	if bid == 0 || itid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	reciever := s.GetParam(r, "reciever").String()
	address := s.GetParam(r, "address").String()
	phone := s.GetParam(r, "phone").String()

	counts := s.GetParam(r, "counts").Int()
	if counts == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_COUNTS)
		return
	}
	additional := s.GetParam(r, "additional").String()

	var (
		name  string
		price int
	)

	if err := s.DB.Query(
		"SELECT name, price FROM item WHERE id = ? LIMIT 1",
		itid,
	).Scan(&name, &price); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	trans, err := s.DB.Begin()
	if err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	var oid int64
	order := &Order{}
	if lastId, err := trans.Execute(
		"INSERT INTO `order` (cuid, bid, itid, name, price, counts, reciever, address, phone, stage, additional) "+
			"VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		cuid,
		bid,
		itid,
		name,
		price,
		counts,
		reciever,
		address,
		phone,
		ORDER_STAGE_NEW,
		additional,
	).LastId(); err != nil {
		trans.Rollback()
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	} else {
		oid = lastId
	}

	if err := trans.Query(
		"SELECT "+strings.Join(order.Properties(), ",")+" "+
			"FROM `order` WHERE id = ? LIMIT 1",
		oid,
	).ScanFunc(func(r Scanner) error {
		return order.Scan(r)
	}); err != nil {
		trans.Rollback()
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	if err := trans.Execute(
		"INSERT INTO order_log (cuid,oid,op,new) "+
			"VALUES (?,?,?,?)",
		cuid,
		oid,
		ORDER_STAGE_NEW,
		order.String(),
	).Error(); err != nil {
		trans.Rollback()
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	if err := trans.Commit(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}
