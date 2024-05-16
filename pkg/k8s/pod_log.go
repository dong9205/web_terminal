package k8s

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
)

func NewWatchContainerLogRequest() *WatchContainerLogRequest {
	return &WatchContainerLogRequest{
		PodLogOptions: &v1.PodLogOptions{
			Follow:                       true,
			Previous:                     false,
			InsecureSkipTLSVerifyBackend: true,
		},
	}
}

// 查看容器日志请求体
type WatchContainerLogRequest struct {
	Namespace string `json:"namespace" validate:"required"`
	PodName   string `json:"pod_name" validate:"required"`
	*v1.PodLogOptions
}

// 查看容器日志
func (c *KubernetesClient) WatchContainerLog(ctx context.Context, req *WatchContainerLogRequest) (io.ReadCloser, error) {
	restReq := c.Client.CoreV1().Pods(req.Namespace).GetLogs(req.PodName, req.PodLogOptions)
	return restReq.Stream(ctx)
}
