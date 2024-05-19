package k8s_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/dong9205/web_terminal/pkg/k8s"
)

func TestServerVersion(t *testing.T) {
	client := setup()
	v, err := client.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func setup() *k8s.KubernetesClient {
	client, err := k8s.NewClientFormFile2("/root/.kube/config", "kind-cluster01")
	if err != nil {
		panic(err)
	}
	return client
}

func TestWatchContainerLog(t *testing.T) {
	client := setup()
	var tailLine int64 = 10
	req := k8s.NewWatchContainerLogRequest()
	req.Namespace = "karmada-system"
	req.PodName = "etcd-0"
	req.TailLines = &tailLine
	stream, err := client.WatchContainerLog(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginContainer(t *testing.T) {
	client := setup()
	reader, write := io.Pipe()

	term := &k8s.MockContainerTerminal{
		In: reader,
	}
	// 模拟用户输入
	go func() {
		write.Write([]byte("ls -al /\n"))
	}()
	req := k8s.NewLoginContainerRequest(term)
	req.Namespace = "karmada-system"
	req.PodName = "karmada-webhook-9cfc4798b-pn8tz"
	err := client.LoginContainer(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNamespaceList(t *testing.T) {
	client := setup()
	req := k8s.NewNamespaceListRequest()
	nss, err := client.NamespaceList(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	for _,ns := range nss.Items{
		t.Log(ns.Name)
	}
}
func TestPodList(t *testing.T) {
	client := setup()
	req := k8s.NewPodListRequest()
	req.Namespace = "default"
	pods, err := client.PodList(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	for _,pod := range pods.Items{
		t.Log(pod.Name)
	}
}
