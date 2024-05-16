package serve

import "net/http"

func Run() {
	http.HandleFunc("/ws/pod/terminal/log", PodTerminalLog)
}
