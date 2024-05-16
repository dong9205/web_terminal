package k8s

import (
	"context"
	"io"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

var shellCmd = []string{
	"/bin/sh",
	"-c",
	`TERM=xterm-256color;export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c "/bin/bash" /dev/null || exec /bin/bash) || exec /bin/sh`,
}

func NewLoginContainerRequest(ce ContainerTerminal) *LoginContainerRequest {
	return &LoginContainerRequest{
		Command:  shellCmd,
		Executor: ce,
	}
}

// 登录容器请求体
type LoginContainerRequest struct {
	Namespace string            `json:"namespace" validate:"required"`
	PodName   string            `json:"pod_name" validate:"required"`
	Container string            `json:"container" validate:"required"`
	Command   []string          `json:"command" validate:"required"`
	Executor  ContainerTerminal `json:"-"`
}

// func (req *LoginContainerRequest) String() string {
// 	return pretty.ToJSON(req)
// }

// 登录容器
func (c *KubernetesClient) LoginContainer(ctx context.Context, req *LoginContainerRequest) error {
	// 构建容器登录请求
	restReq := c.Client.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(req.PodName).
		Namespace(req.Namespace).
		SubResource("exec")
	// 登录容器参数
	restReq.VersionedParams(&v1.PodExecOptions{
		Container: req.Container,
		Command:   req.Command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)
	// 执行容器终端登录
	executor, err := remotecommand.NewSPDYExecutor(c.Restconf, "POST", restReq.URL())
	if err != nil {
		return err
	}
	return executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             req.Executor,
		Stdout:            req.Executor,
		Stderr:            req.Executor,
		Tty:               true,
		TerminalSizeQueue: req.Executor,
	})
}

type ContainerTerminal interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type MockContainerTerminal struct {
	In io.Reader
}

func (t *MockContainerTerminal) Read(p []byte) (n int, err error) {
	return t.In.Read(p)
}

func (t *MockContainerTerminal) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (t *MockContainerTerminal) Next() *remotecommand.TerminalSize {
	return &remotecommand.TerminalSize{
		Width:  100,
		Height: 100,
	}
}
