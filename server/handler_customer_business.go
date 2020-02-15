package server

import (
	"net/http"
	"strings"
)

func (s *Server) CustomerBusinessGet(w http.ResponseWriter, r *http.Request) {
	pager := s.GetPager(r)

	list := []*Business{}
	var business *Business
	if err := s.DB.Query(
		"SELECT " + strings.Join(business.Properties(), ",") + " " +
			"FROM `business_profile` WHERE status = 1" + pager.Sql(),
	).ScanFunc(func(r Scanner) error {
		business = &Business{}
		if err := business.Scan(r); err != nil {
			return err
		}
		list = append(list, business)
		return nil
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, list)
}

func (s *Server) CustomerBusinessBidGet(w http.ResponseWriter, r *http.Request) {
	bid := s.GetParam(r, "bid").Int()
	if bid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	business := &Business{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(business.Properties(), ",")+" "+
			"FROM business_profile WHERE id = ? AND status = 1 LIMIT 1",
		bid,
	).ScanFunc(func(r Scanner) error {
		return business.Scan(r)
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, business)
}
