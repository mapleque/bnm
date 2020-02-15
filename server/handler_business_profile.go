package server

import (
	"net/http"
	"strings"
)

func (s *Server) BusinessProfileGet(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()

	business := &Business{}
	if err := s.DB.Query(
		"SELECT "+strings.Join(business.Properties(), ", ")+" "+
			"FROM business_profile WHERE buid = ? LIMIT 1",
		buid,
	).ScanFunc(func(r Scanner) error {
		return business.Scan(r)
	}); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}
	s.Success(w, business)
}

func (s *Server) BusinessProfilePost(w http.ResponseWriter, r *http.Request) {
	buid := s.GetParam(r, "buid").Int()

	name := s.GetParam(r, "name").String()
	avatar := s.GetParam(r, "avatar").String()
	desc := s.GetParam(r, "desc").String()
	qrcode := s.GetParam(r, "qrcode").String()
	wxid := s.GetParam(r, "wxid").String()
	if name == "" || (qrcode == "" && wxid == "") {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_BUSINESS_PROFILE)
		return
	}

	if exist, err := s.DB.Query(
		"SELECT * FROM business_profile WHERE buid = ? LIMIT 1",
		buid,
	).Exist(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	} else if !exist {
		if err := s.DB.Execute(
			"INSERT INTO business_profile "+
				"(buid, name, avatar, `desc`, qrcode, wxid) "+
				"VALUES (?,?,?,?,?,?)",
			buid,
			name,
			avatar,
			desc,
			qrcode,
			wxid,
		).Error(); err != nil {
			s.Response(w, STATUS_SERVER_ERROR, err)
			return
		}
	} else {
		if err := s.DB.Execute(
			"UPDATE business_profile "+
				"SET avatar = ?, `desc` = ? "+
				"WHERE buid = ?",
			avatar,
			desc,
			buid,
		).Error(); err != nil {
			s.Response(w, STATUS_SERVER_ERROR, err)
			return
		}
	}

	s.Success(w, nil)
}
