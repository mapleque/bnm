package server

func (s *Server) initRouter() {
	s.handle("/business/login", "POST", s.BusinessLogin)
	s.handle("/business/profile", "GET", s.BusinessAuth, s.BusinessProfileGet)
	s.handle("/business/profile", "POST", s.BusinessAuth, s.FormParam, s.BusinessProfilePost)
	s.handle("/business/items", "GET", s.BusinessAuth, s.PagerParam, s.BusinessItemsGet)
	s.handle("/business/items", "POST", s.BusinessAuth, s.FormParam, s.BusinessItemsPost)
	s.handle("/business/items/<itid>", "GET", s.BusinessAuth, s.BusinessItemsItidGet)
	s.handle("/business/items/<itid>", "POST", s.BusinessAuth, s.FormParam, s.BusinessItemsItidPost)
	s.handle("/business/orders", "GET", s.BusinessAuth, s.PagerParam, s.BusinessOrdersGet)
	s.handle("/business/orders/<oid>", "GET", s.BusinessAuth, s.BusinessOrdersOidGet)
	s.handle("/business/orders/<oid>", "POST", s.BusinessAuth, s.FormParam, s.BusinessOrdersOidPost)

	s.handle("/customer/login", "POST", s.CustomerLogin)
	s.handle("/customer/business", "GET", s.CustomerAuth, s.PagerParam, s.CustomerBusinessGet)
	s.handle("/customer/business/<bid>", "GET", s.CustomerAuth, s.CustomerBusinessBidGet)
	s.handle("/customer/business/<bid>/items", "GET", s.CustomerAuth, s.PagerParam, s.CustomerBusinessBidItemsGet)
	s.handle("/customer/business/<bid>/items/<itid>", "GET", s.CustomerAuth, s.CustomerBusinessBidItemsItidGet)
	s.handle("/customer/business/<bid>/items/<itid>", "POST", s.CustomerAuth, s.FormParam, s.CustomerBusinessBidItemsItidPost)
	s.handle("/customer/orders", "GET", s.CustomerAuth, s.PagerParam, s.CustomerOrdersGet)
	s.handle("/customer/orders/<oid>", "GET", s.CustomerAuth, s.CustomerOrdersOidGet)
	s.handle("/customer/orders/<oid>", "POST", s.CustomerAuth, s.FormParam, s.CustomerOrdersOidPost)
	s.handle("/customer/address", "GET", s.CustomerAuth, s.PagerParam, s.CustomerAddressGet)
	s.handle("/customer/address", "POST", s.CustomerAuth, s.FormParam, s.CustomerAddressPost)
	s.handle("/customer/address/<aid>", "POST", s.CustomerAuth, s.FormParam, s.CustomerAddressAidPost)
	s.handle("/customer/address/<aid>", "DELETE", s.CustomerAuth, s.CustomerAddressAidDelete)

}
