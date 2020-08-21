package server

import (
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	code := s.GetParam(r, "code").String()
	if code == "" {
		s.Response(w, STATUS_SERVER_ERROR, MSG_INVALID_CODE)
		return
	}

	// request open_id from wx api
	auth, err := s.Wx.Jscode2session(code)
	if err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	// randome a token as session key
	token := GetSessionToken()

	if exist, err := s.DB.Query(
		"SELECT * FROM user WHERE open_id = ? LIMIT 1",
		auth.OpenId,
	).Exist(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	} else if !exist {
		// create a new user data
		if err := s.DB.Execute(
			"INSERT INTO user "+
				"(open_id, union_id, token, expired_at) "+
				"VALUES (?,?,?,DATE_ADD(NOW(),INTERVAL 7 DAY))",
			auth.OpenId,
			auth.UnionId,
			token,
		).Error(); err != nil {
			s.Response(w, STATUS_SERVER_ERROR, err)
			return
		}
	} else {
		if err := s.DB.Execute(
			"UPDATE user "+
				"SET token = ?, expired_at = DATE_ADD(NOW(),INTERVAL 7 DAY) "+
				"WHERE open_id = ? LIMIT 1",
			token,
			auth.OpenId,
		).Error(); err != nil {
			s.Response(w, STATUS_SERVER_ERROR, err)
			return
		}
	}

	s.Success(w, token)
}

func (s *Server) Auth(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Session-Key")
	if token == "" {
		s.Response(w, STATUS_FORBIDDEN, MSG_FORBIDDEN)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}

	var buid int
	if err := s.DB.Query(
		"SELECT id FROM user WHERE token = ? AND expired_at > NOW() LIMIT 1",
		token,
	).Scan(&buid); err != nil {
		s.Response(w, STATUS_FORBIDDEN, MSG_LOGIN_EXPIRED)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}

	s.SaveParam(r, "buid", buid)
	s.SaveParam(r, "cuid", buid)
}
