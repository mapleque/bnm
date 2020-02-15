package server

const (
	// show for user
	MSG_INVALID_BUSINESS_PROFILE = "请正确填写认证信息"
	MSG_INVALID_ITEM             = "请填写完整的商品信息"
	MSG_FORBIDDEN                = "请登陆授权后访问"
	MSG_LOGIN_EXPIRED            = "登陆已过期，请重新登陆"
	MSG_INVALID_COUNTS           = "请填写正确的购买商品数量"
	MSG_INVALID_ADDRESS          = "请填写正确的商品邮寄地址"

	// show for hacker
	MSG_INVALID_CODE  = "code不合法"
	MSG_INVALID_PAGER = "错误的分页参数"
	MSG_INVALID_BODY  = "请求体格式不正确"
	MSG_INVALID_ID    = "错误的资源标识"
	MSG_INVALID_OP    = "非法的操作"
)
