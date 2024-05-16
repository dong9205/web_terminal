package serve

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Request struct {
	// 请求ID
	ID string `json:"id"`
	// 指令名称
	Command string `json:"command"`
	// 指令参数
	Params string `json:"params"`
}

func PodTerminalLog(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	defer r.Body.Close()
	requestStr, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var req Request
	err = json.Unmarshal(requestStr, &req)
	if err != nil {
		log.Println("unmarshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
