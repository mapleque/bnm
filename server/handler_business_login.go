package server

import (
	"net/http"
)

func (s *Server) BusinessLogin(w http.ResponseWriter, r *http.Request) {
	code := s.GetParam(r, "code").String()
	if code == "" {
		s.Response(w, STATUS_SERVER_ERROR, MSG_INVALID_CODE)
		return
	}

	// request open_id from wx api
	auth, err := s.Wx2B.Jscode2session(code)
	if err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	}

	// randome a token as session key
	token := GetSessionToken()

	if exist, err := s.DB.Query(
		"SELECT * FROM business_user WHERE open_id = ? LIMIT 1",
		auth.OpenId,
	).Exist(); err != nil {
		s.Response(w, STATUS_SERVER_ERROR, err)
		return
	} else if !exist {
		// create a new user data
		if err := s.DB.Execute(
			"INSERT INTO business_user "+
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
			"UPDATE business_user "+
				"SET token = ?, expired_at = DATE_ADD(NOW(),INTERVAL 7 DAY) "+
				"WHERE openid = ? LIMIT 1",
			token,
			auth.OpenId,
		).Error(); err != nil {
			s.Response(w, STATUS_SERVER_ERROR, err)
			return
		}
	}

	s.Success(w, token)
}

func (s *Server) BusinessAuth(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("SessionKey")
	if token == "" {
		s.Response(w, STATUS_FORBIDDEN, MSG_FORBIDDEN)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}

	var buid int
	if err := s.DB.Query(
		"SELECT id FROM business_user WHERE token = ? AND expired_at > NOW() LIMIT 1",
		token,
	).Scan(&buid); err != nil {
		s.Response(w, STATUS_FORBIDDEN, MSG_LOGIN_EXPIRED)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}

	s.SaveParam(r, "buid", buid)
}
