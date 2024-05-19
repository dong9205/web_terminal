package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodListRequest struct {
	Namespace string `json:"namespace" validate:"required"`
	metav1.ListOptions
}

func NewPodListRequest() *PodListRequest {
	return &PodListRequest{}
}

func (c *KubernetesClient) PodList(ctx context.Context,podListRequest *PodListRequest) (*v1.PodList,error) {
	return c.Client.CoreV1().Pods(podListRequest.Namespace).List(ctx, podListRequest.ListOptions)
}