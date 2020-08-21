package server

import (
	"net/http"
	"strings"
)

func (s *Server) BusinessOrdersGet(w http.ResponseWriter, r *http.Request) {
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

	list := []*Order{}
	var order *Order
	if err := s.DB.Query(
		"SELECT "+strings.Join(order.Properties(), ",")+" "+
			"FROM `order` WHERE bid = ? "+pager.Sql(),
		bid,
	).ScanFunc(func(r Scanner) error {
		order = &Order{}
		if err := order.Scan(r); err != nil {
			return err
		}
		list = append(list, order)
		return nil
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, list)
}

func (s *Server) BusinessOrdersOidGet(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()
	oid := s.GetParam(r, "oid").Int()
	if oid == 0 {
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

	order := &Order{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(order.Properties(), ", ")+" "+
			"FROM `order` WHERE bid = ? AND id = ? LIMIT 1",
		bid,
		oid,
	).ScanFunc(func(r Scanner) error {
		return order.Scan(r)
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, order)
}

func (s *Server) BusinessOrdersOidPost(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()
	oid := s.GetParam(r, "oid").Int()
	if oid == 0 {
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

	price := s.GetParam(r, "price").Int()
	expNo := s.GetParam(r, "exp_no").String()
	additional := s.GetParam(r, "additional").String()
	reciever := s.GetParam(r, "reciever").String()
	address := s.GetParam(r, "address").String()
	phone := s.GetParam(r, "phone").String()

	stage := s.GetParam(r, "stage").String()
	var (
		status   int
		preStage []interface{}
	)
	switch stage {
	case ORDER_STAGE_NEW, ORDER_STAGE_PAID:
		preStage = []interface{}{ORDER_STAGE_NEW}
		status = 0
	case ORDER_STAGE_B_CANCEL:
		preStage = []interface{}{ORDER_STAGE_NEW}
		status = 1
	case ORDER_STAGE_TRANSPORT:
		preStage = []interface{}{
			ORDER_STAGE_TRANSPORT,
			ORDER_STAGE_PAID,
			ORDER_STAGE_C_W_CANCEL,
		}
		status = 0
	case ORDER_STAGE_REPAID:
		preStage = []interface{}{
			ORDER_STAGE_PAID,
			ORDER_STAGE_C_W_CANCEL,
		}
		status = 1
	default:
		s.Response(w, STATUS_SERVER_ERROR, MSG_INVALID_OP)
		return
	}

	trans, err := s.DB.Begin()
	if err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	order := &Order{}
	bind := []interface{}{
		price,
		expNo,
		additional,
		reciever,
		address,
		phone,
		stage,
		status,
		oid,
		bid,
	}
	bind = append(bind, preStage...)
	if eff, err := trans.Execute(
		"UPDATE `order` SET price=?,exp_no=?,additional=?,reciever=?,address=?,phone=?,stage=?,status=? "+
			"WHERE id=? AND bid=? AND stage IN ("+
			strings.Repeat(",?", len(preStage))[1:]+
			") LIMIT 1",
		bind...,
	).EffectRows(); err != nil {
		trans.Rollback()
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	} else if eff != 1 {
		trans.Rollback()
		s.Response(w, STATUS_SERVER_ERROR, MSG_INVALID_OP)
		return
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
		"INSERT INTO order_log (buid,oid,op,new) "+
			"VALUES (?,?,?,?)",
		buid,
		oid,
		stage,
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
