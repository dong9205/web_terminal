package terminal

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Request struct {
	// 请求ID
	ID string `json:"id"`
	// 指令名称
	Command string `json:"command"`
	// 指令参数
	Params json.RawMessage `json:"params"`
}

type Response struct {
	Request *Request `json:"request"`
	// 异常信息
	Message string `json:"message"`
	// 处理完成的数据
	Data any `json:"data"`
}

func NewResponse() *Response {
	return &Response{}
}

type HandleFunc func(*Request, *Response)

var handlerFuncs = map[string]HandleFunc{}

// 注册请求处理函数
func RegisterCmdHandlerFunc(command string, fn HandleFunc) {
	handlerFuncs[command] = fn
}

func GetHandlerFunc(command string) HandleFunc {
	return handlerFuncs[command]
}

type WebSocketTerminal struct {
	ws *websocket.Conn
	*TerminalResizer

	// Write所需属性
	timeout  time.Duration
	writeBuf []byte
}

func NewWebSocketTerminal(ws *websocket.Conn) *WebSocketTerminal {
	return &WebSocketTerminal{
		ws:              ws,
		TerminalResizer: NewTerminalSizer(),
		timeout:         time.Second * 3,
		writeBuf:        make([]byte, DefaultWriteBuf),
	}
}

func (t *WebSocketTerminal) Close() {
	t.ws.Close()
}

func (t *WebSocketTerminal) Read(p []byte) (n int, err error) {
	mt, m, err := t.ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	switch mt {
	case websocket.TextMessage:
		t.HandleCmd(m)
	case websocket.CloseMessage:
		log.Printf(`receive client close: %s\n`, m)
	default:
		n = copy(p, m)
	}
	return n, nil
}

func (t *WebSocketTerminal) Response(resp *Response) {
	t.ws.WriteMessage(websocket.TextMessage, []byte(resp.Message))
}

func (t *WebSocketTerminal) HandleCmd(m []byte) {
	resp := NewResponse()
	defer t.Response(resp)

	req, err := ParseRequest(m)
	if err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Request = req

	switch req.Command {
	case "resize":
		payload := NewTerminalSize()
		err := json.Unmarshal(req.Params, payload)
		if err != nil {
			resp.Message = err.Error()
			return
		}
		t.SetSize(*payload)
		return
	}
	// 处理自定义命令
	if fn := GetHandlerFunc(req.Command); fn != nil {
		fn(req, resp)
		return
	}
	resp.Message = "command not found"
}

func ParseRequest(m []byte) (*Request, error) {
	var req Request
	err := json.Unmarshal(m, &req)
	return &req, err
}
