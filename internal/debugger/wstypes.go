package debugger

const WsMessageTypeDebugger WsMessageType = 1

type WsMessageType int

type WsMessagePayload interface{ WsDebuggerMessage }
type WsMessage[P WsMessagePayload] struct {
	Type    WsMessageType `json:"type,omitempty"`
	Payload P             `json:"payload,omitempty"`
}

type WsDebuggerMessage struct {
	Type     WsDebuggerMessageType `json:"type,omitempty"`
	Imposter *string               `json:"imposter,omitempty"`
	Request  *string               `json:"request,omitempty"`
	Response *string               `json:"response,omitempty"`
}

const (
	WsDebuggerMessageTypeRequestReceived WsDebuggerMessageType = iota + 1
	WsDebuggerMessageTypeRequestConfirmed
	WsDebuggerMessageTypeImposterMatched
	WsDebuggerMessageTypeImposterConfirmed
	WsDebuggerMessageTypeResponsePrepared
	WsDebuggerMessageTypeResponseConfirmed
)

type WsDebuggerMessageType int
