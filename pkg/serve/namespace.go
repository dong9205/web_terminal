package serve

import (
	"net/http"

	"github.com/dong9205/web_terminal/pkg/k8s"
)

func NamespaceList(w http.ResponseWriter, r *http.Request) {
	// 读取请求参数解析到结构体
	req := k8s.NewNamespaceListRequest()
	if err := r.ParseForm(); err != nil {
		response(w, http.StatusBadRequest, nil, err)
		return
	}
	// 获取Kubernetes客户端
	k8sClient, err := k8s.NewClientFormFile2("/root/.kube/config", "kind-cluster01")
	if err != nil {
		response(w, http.StatusInternalServerError, nil, err)
		return
	}
	nss, err := k8sClient.NamespaceList(r.Context(), req)
	if err != nil {
		response(w, http.StatusInternalServerError, nil, err)
		return
	}
	var nssResp []string = make([]string, 0)
	for _, ns := range nss.Items {
		nssResp = append(nssResp, ns.Name)
	}
	response(w, http.StatusOK, nssResp, nil)
}
