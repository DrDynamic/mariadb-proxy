package tcp

type ProxyRecEvent struct {
	Data []byte
}

func NewProxyRecEvent(data []byte) *ProxyRecEvent {
	return &ProxyRecEvent{
		Data: data,
	}
}
