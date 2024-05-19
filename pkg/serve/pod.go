package serve

import (
	"context"
	"errors"
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

func PodList(w http.ResponseWriter, r *http.Request) {
	// 读取请求参数解析到结构体
	req := k8s.NewPodListRequest()
	if err := r.ParseForm(); err != nil {
		response(w, http.StatusBadRequest, nil, err)
		return
	}
	req.Namespace = r.Form.Get("namespace")
	if req.Namespace == "" {
		response(w, http.StatusBadRequest, nil, errors.New("namespace is required"))
		return
	}
	// 获取Kubernetes客户端
	k8sClient, err := k8s.NewClientFormFile2("/root/.kube/config", "kind-cluster01")
	if err != nil {
		response(w, http.StatusInternalServerError, nil, err)
		return
	}
	pods, err := k8sClient.PodList(r.Context(), req)
	if err != nil {
		response(w, http.StatusInternalServerError, nil, err)
		return
	}
	var podListResp map[string][]string = make(map[string][]string)
	for _, pod := range pods.Items {
		podList := make([]string, 0)
		for _, container := range pod.Spec.Containers {
			podList = append(podList, container.Name)
		}
		podListResp[pod.Name] = podList
	}
	response(w, http.StatusOK, podListResp, nil)
}
