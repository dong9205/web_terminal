package terminal

func init() {
	RegisterCmdHandlerFunc("ping", PingHandlerFunc)
}

func PingHandlerFunc(r *Request, w *Response) {
	w.Data = "pong"
}
