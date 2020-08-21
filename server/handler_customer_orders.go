package server

import (
	"net/http"
	"strings"
)

func (s *Server) CustomerOrdersGet(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	pager := s.GetPager(r)

	list := []*Order{}
	var order *Order
	if err := s.DB.Query(
		"SELECT "+strings.Join(order.Properties(), ",")+" "+
			"FROM `order` WHERE cuid = ? "+pager.Sql(),
		cuid,
	).ScanFunc(func(r Scanner) error {
		order := &Order{}
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

func (s *Server) CustomerOrdersOidGet(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	oid := s.GetParam(r, "oid").Int()
	if oid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}

	order := &Order{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(order.Properties(), ", ")+" "+
			"FROM `order` WHERE cuid = ? AND id = ? LIMIT 1",
		cuid,
		oid,
	).ScanFunc(func(r Scanner) error {
		return order.Scan(r)
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, order)
}

func (s *Server) CustomerOrdersOidPost(w http.ResponseWriter, r *http.Request) {
	cuid := s.GetParam(r, "cuid").Int()
	oid := s.GetParam(r, "oid").Int()
	if oid == 0 {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_ID)
		return
	}
	stage := s.GetParam(r, "stage").String()
	var (
		status   int
		preStage []interface{}
	)
	switch stage {
	case ORDER_STAGE_C_CANCEL:
		preStage = []interface{}{
			ORDER_STAGE_NEW,
		}
		status = 1
	case ORDER_STAGE_C_W_CANCEL:
		preStage = []interface{}{
			ORDER_STAGE_PAID,
		}
		status = 0
	case ORDER_STAGE_FINISH:
		preStage = []interface{}{
			ORDER_STAGE_TRANSPORT,
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
		stage,
		status,
		oid,
		cuid,
	}
	bind = append(bind, preStage...)
	if eff, err := trans.Execute(
		"UPDATE `order` SET stage=?,status=? "+
			"WHERE id=? AND cuid=? AND stage IN ("+
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
		"INSERT INTO order_log (cuid,oid,op,new) "+
			"VALUES (?,?,?,?)",
		cuid,
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
