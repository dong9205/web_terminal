package serve

import (
	"encoding/json"
	"net/http"
)

func Run() error {
	http.HandleFunc("/api/namespace/list", NamespaceList)
	http.HandleFunc("/api/pod/list", PodList)
	http.HandleFunc("/ws/pod/terminal/log", PodTerminalLog)
	http.HandleFunc("/ws/pod/terminal/login", PodTerminalLogin)
	http.Handle("/", http.FileServer(http.FS(getUiWeb())))
	return http.ListenAndServe(":9200", nil)
}

type responseStruct struct {
	Code int    `json:"code"`
	Data any    `json:"data,omitempty"`
	Err  string `json:"err,omitempty"`
}

func response(w http.ResponseWriter, code int, data any, err error) {
	w.WriteHeader(http.StatusOK)
	var resp responseStruct
	if err != nil {
		resp = responseStruct{
			Code: code,
			Err:  err.Error(),
		}
	} else {
		resp = responseStruct{
			Code: code,
			Data: data,
		}
	}
	responseByte, _ := json.Marshal(resp)
	w.Write(responseByte)
}
