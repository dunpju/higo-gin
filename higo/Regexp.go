package higo

// 正则常量
const (
	// 开始符号
	BEGIN_SYMBOL = "^"
	// 结束符号
	END_SYMBOL = "$"
	// 冒号
	COLON = ":"
	// IP 地址
	IP = "(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)"
	// 合法端口【0-65535】
	PORT = "([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])"
	// IP:PORT
	IP_PORT = BEGIN_SYMBOL + IP + COLON + PORT + END_SYMBOL
	// :PORT
	COLON_PORT = BEGIN_SYMBOL + COLON + PORT + END_SYMBOL
)
