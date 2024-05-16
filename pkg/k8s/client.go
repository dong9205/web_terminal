package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func NewClientFormFile2(configPath, useContext string) (*KubernetesClient, error) {
	rawConfig, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}
	if useContext != "" {
		rawConfig.CurrentContext = useContext
	}
	restConf, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*api.Config, error) {
		return rawConfig, nil
	})
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}
	return &KubernetesClient{
		Client:   client,
		Kubeconf: rawConfig,
		Restconf: restConf,
	}, nil
}

func NewClientFormFile(configPath string) (*KubernetesClient, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	kc, err := os.ReadFile(filepath.Join(wd, configPath))
	if err != nil {
		return nil, err
	}
	return NewClient(kc)
}

func NewClient(kubeConfigContent []byte) (*KubernetesClient, error) {
	// 加载Kuberconfig配置
	kubeConf, err := clientcmd.Load(kubeConfigContent)
	if err != nil {
		return nil, err
	}
	// 构造Resetclinet Config
	resetConf, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*api.Config, error) {
		return kubeConf, nil
	})
	if err != nil {
		return nil, err
	}

	// 初始化客户端
	client, err := kubernetes.NewForConfig(resetConf)
	if err == nil {
		return nil, err
	}
	return &KubernetesClient{
		Client:   client,
		Kubeconf: kubeConf,
		Restconf: resetConf,
	}, nil
}

type KubernetesClient struct {
	Kubeconf *api.Config
	Restconf *rest.Config
	Client   *kubernetes.Clientset
}

// 查看当前Service端版本
func (c *KubernetesClient) ServerVersion() (string, error) {
	si, err := c.Client.ServerVersion()
	if err != nil {
		return "", err
	}

	return si.String(), nil
}
