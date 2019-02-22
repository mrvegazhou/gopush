package constdefine

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	IM_UNAUTHORIZED: "unauthorized", 	//需要认证
	IM_BAD_REQUEST: "bad request", 		//请求错误
	IM_UNKNOWN_ERROR:	"unkown error",	//未知错误
	IM_DEVICE_NOT_BIND_USER:	"device not bind user",	// 设备没有绑定用户
	IM_ERROR_DEVICE_TOKEN:	"error device token",// 设备id或者token错误
	NUMBER_HAS_BE_USED:	"number has be used",// 手机号码已经被使用
	ERROR_NUMBER_OR_PASSWORD:	"error number or password",// 用户名或者密码错误
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}