package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceListRequest struct{
	metav1.ListOptions
}

func NewNamespaceListRequest() *NamespaceListRequest {
	return &NamespaceListRequest{}
}

func (c *KubernetesClient)NamespaceList(ctx context.Context,namespaceListRequest *NamespaceListRequest)(*v1.NamespaceList, error) {
	return c.Client.CoreV1().Namespaces().List(ctx, namespaceListRequest.ListOptions)
}