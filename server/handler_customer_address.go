package server

import (
	"net/http"
	"strings"
)

func (s *Server) CustomerAddressGet(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	pager := s.GetPager(r)

	list := []*Address{}
	var address *Address
	if err := s.DB.Query(
		"SELECT "+strings.Join(address.Properties(), ",")+" "+
			"FROM customer_address WHERE cuid= ? "+pager.Sql(),
		cuid,
	).ScanFunc(func(r Scanner) error {
		address = &Address{}
		if err := address.Scan(r); err != nil {
			return err
		}
		list = append(list, address)
		return nil
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, list)
}

func (s *Server) CustomerAddressPost(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	label := s.GetParam(r, "label").String()
	reciever := s.GetParam(r, "reciever").String()
	address := s.GetParam(r, "address").String()
	phone := s.GetParam(r, "phone").String()

	if label == "" {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ADDRESS)
		return
	}
	if err := s.DB.Execute(
		"INSERT INTO customer_address (cuid, label, reciever, address, phone) "+
			"VALUES (?,?,?,?,?)",
		cuid,
		label,
		reciever,
		address,
		phone,
	).Error(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}

func (s *Server) CustomerAddressAidPost(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	aid := s.GetParam(r, "aid").Int()
	if aid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}
	label := s.GetParam(r, "label").String()
	reciever := s.GetParam(r, "reciever").String()
	address := s.GetParam(r, "address").String()
	phone := s.GetParam(r, "phone").String()
	if label == "" {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ADDRESS)
		return
	}
	if err := s.DB.Execute(
		"UPDATE customer_address SET label=?,reciever=?,address=?,phone=? WHERE id=? AND cuid=? LIMIT 1",
		label,
		reciever,
		address,
		phone,
		aid,
		cuid,
	).Error(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}

func (s *Server) CustomerAddressAidDelete(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	aid := s.GetParam(r, "aid").Int()
	if aid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}
	if err := s.DB.Execute(
		"DELETE FROM customer_address WHERE id=? AND cuid=? LIMIT 1",
		aid,
		cuid,
	).Error(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, nil)
}
