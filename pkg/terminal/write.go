package terminal

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

var DefaultWriteBuf = 4 * 1024

func (i *WebSocketTerminal) ReadReq(req any) error {
	mt, data, err := i.ws.ReadMessage()
	if err != nil {
		return err
	}
	if mt != websocket.TextMessage {
		return fmt.Errorf(`req must be TextMessage, but now not ,is %d`, mt)
	}
	if !json.Valid(data) {
		return fmt.Errorf(`req must be json data, but %s`, string(data))
	}
	return json.Unmarshal(data, req)
}

// func (i *WebSocketTerminal) WriteTo(req any) error {
// }

func (i *WebSocketTerminal) Write(p []byte) (n int, err error) {
	err = i.ws.WriteMessage(websocket.BinaryMessage, p)
	n = len(p)
	return
}

func (i *WebSocketTerminal) Failed(err error) {
	i.ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	i.Close()
}
