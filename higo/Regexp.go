package higo

// 正则常量
const (
	// BEGIN_SYMBOL 开始符号
	BEGIN_SYMBOL = "^"
	// END_SYMBOL 结束符号
	END_SYMBOL = "$"
	// COLON 冒号
	COLON = ":"
	// IP 地址
	IP = "(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)"
	// PORT 合法端口【0-65535】
	PORT = "([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])"
	// IP_PORT IP:PORT
	IP_PORT = BEGIN_SYMBOL + IP + COLON + PORT + END_SYMBOL
	// COLON_PORT :PORT
	COLON_PORT = BEGIN_SYMBOL + COLON + PORT + END_SYMBOL
)
