package serve

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dong9205/web_terminal/pkg/k8s"
	"github.com/dong9205/web_terminal/pkg/terminal"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 60,
	ReadBufferSize:   8192,
	WriteBufferSize:  8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PodTerminalLog(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	term := terminal.NewWebSocketTerminal(ws)

	// 读取前端传的参数，如果有认证也可以加在这里
	req := k8s.NewWatchContainerLogRequest()
	if err := term.ReadReq(req); err != nil {
		term.Failed(err)
	}
	log.Printf(`watch log req %v\n`, req)
	// 获取Kubernetes客户端
	k8sClient, err := k8s.NewClientFormFile2("/root/.kube/config", "kind-cluster01")
	if err != nil {
		term.Failed(err)
	}
	stream, err := k8sClient.WatchContainerLog(context.Background(), req)
	if err != nil {
		term.Failed(err)
	}
	_, err = io.Copy(term, stream)
	if err != nil {
		term.Failed(err)
	}
}

func PodTerminalLogin(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	term := terminal.NewWebSocketTerminal(ws)

	// 读取前端传的参数，如果有认证也可以加在这里
	req := k8s.NewLoginContainerRequest(term)
	if err := term.ReadReq(req); err != nil {
		term.Failed(err)
	}
	log.Printf(`watch login req %v\n`, req)
	// 获取Kubernetes客户端
	k8sClient, err := k8s.NewClientFormFile2("/root/.kube/config", "kind-cluster01")
	if err != nil {
		term.Failed(err)
	}
	err = k8sClient.LoginContainer(context.Background(), req)
	if err != nil {
		term.Failed(err)
	}
}
