package server

func (s *Server) initRouter() {
	s.handle("/login", "POST", s.FormParam, s.Login)

	s.handle("/business/profile", "GET", s.Auth, s.BusinessProfileGet)
	s.handle("/business/profile", "POST", s.Auth, s.FormParam, s.BusinessProfilePost)
	s.handle("/business/items", "GET", s.Auth, s.PagerParam, s.BusinessItemsGet)
	s.handle("/business/items", "POST", s.Auth, s.FormParam, s.BusinessItemsPost)
	s.handle("/business/items/<itid>", "GET", s.Auth, s.BusinessItemsItidGet)
	s.handle("/business/items/<itid>", "POST", s.Auth, s.FormParam, s.BusinessItemsItidPost)
	s.handle("/business/orders", "GET", s.Auth, s.PagerParam, s.BusinessOrdersGet)
	s.handle("/business/orders/<oid>", "GET", s.Auth, s.BusinessOrdersOidGet)
	s.handle("/business/orders/<oid>", "POST", s.Auth, s.FormParam, s.BusinessOrdersOidPost)

	s.handle("/customer/business", "GET", s.Auth, s.PagerParam, s.CustomerBusinessGet)
	s.handle("/customer/business/<bid>", "GET", s.Auth, s.CustomerBusinessBidGet)
	s.handle("/customer/business/<bid>/items", "GET", s.Auth, s.PagerParam, s.CustomerBusinessBidItemsGet)
	s.handle("/customer/business/<bid>/items/<itid>", "GET", s.Auth, s.CustomerBusinessBidItemsItidGet)
	s.handle("/customer/business/<bid>/items/<itid>", "POST", s.Auth, s.FormParam, s.CustomerBusinessBidItemsItidPost)
	s.handle("/customer/orders", "GET", s.Auth, s.PagerParam, s.CustomerOrdersGet)
	s.handle("/customer/orders/<oid>", "GET", s.Auth, s.CustomerOrdersOidGet)
	s.handle("/customer/orders/<oid>", "POST", s.Auth, s.FormParam, s.CustomerOrdersOidPost)
	s.handle("/customer/address", "GET", s.Auth, s.PagerParam, s.CustomerAddressGet)
	s.handle("/customer/address", "POST", s.Auth, s.FormParam, s.CustomerAddressPost)
	s.handle("/customer/address/<aid>", "POST", s.Auth, s.FormParam, s.CustomerAddressAidPost)
	s.handle("/customer/address/<aid>", "DELETE", s.Auth, s.CustomerAddressAidDelete)
}
