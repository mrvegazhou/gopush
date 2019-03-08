package constdefine

import "time"

//连接类型
const (
	ConnType_WebSocket = "websocket"
	ConnType_LongPoll  = "longpoll"

	KeyDeviceId = "device_id"
	KeyUserId   = "user_id"

	IMTypeLen		  = 2                 // 消息类型字节数组长度
	IMLenLen        = 2                 // 消息长度字节数组长度
	IMHeadLen       = 4                 // 消息头部字节数组长度（消息类型字节数组长度+消息长度字节数组长度）
	IMContentMaxLen = 4092              // 消息体最大长度
	IMBufLen        = IMContentMaxLen + 4 // 缓冲buffer字节数组长度
	IMReadDeadline  = 10 * time.Minute
	IMWriteDeadline = 10 * time.Second

	IMCodeSignIn         = 1 // 设备登录
	IMCodeSignInACK      = 2 // 设备登录回执
	IMCodeSyncTrigger    = 3 // 消息同步触发
	IMCodeHeadbeat       = 4 // 心跳
	IMCodeHeadbeatACK    = 5 // 心跳回执
	IMCodeMessageSend    = 6 // 消息发送
	IMCodeMessageSendACK = 7 // 消息发送回执
	IMCodeMessage        = 8 // 消息投递
	IMCodeMessageACK     = 9 // 消息投递回执
)

const (
	no int32 = iota
	yes
)
