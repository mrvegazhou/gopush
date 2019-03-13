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
	IM_ERROR_USER_REGIST:	"error user regist",//	用户注册失败
	IM_ERROR_OUT_OF_SIZE:	"package content out of size", // package的content字节数组过大
	NUMBER_HAS_BE_USED:	"number has be used",// 手机号码已经被使用
	ERROR_MOBILE:	"error mobile",// 用户名手机号错误
	ERROR_PASSWORD:	"error password", //用户名密码错误
	ERROR_CONFIG:	"error config", //没有找到配置信息
	ERROR_REDIS: "error redis connection", //redis连接失败

	RECORD_NOT_FOUND:"record not found", //没有数据
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}