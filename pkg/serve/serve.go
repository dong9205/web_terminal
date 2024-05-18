package serve

import "net/http"

func Run() error {
	http.HandleFunc("/ws/pod/terminal/log", PodTerminalLog)
	http.HandleFunc("/ws/pod/terminal/login", PodTerminalLogin)
	http.Handle("/", http.FileServer(http.FS(getUiWeb())))
	return http.ListenAndServe(":9200", nil)
}
