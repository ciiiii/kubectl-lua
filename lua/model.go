package lua

import (
	"flag"
	"os"

	"github.com/ciiiii/kubectl-lua/api"
	lua "github.com/yuin/gopher-lua"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "absolute path of the kubeconfig file")
	flag.Parse()
}

type KubeClient struct {
	clientset         *kubernetes.Clientset
	dynamic           dynamic.Interface
	resourceDiscovery *api.ResourceDiscovery
}

func kubeClientFromConfig() (*KubeClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	resourceDiscovery, err := api.NewResourceDiscovery(clientset)
	if err != nil {
		return nil, err
	}
	return &KubeClient{
		clientset:         clientset,
		dynamic:           dynamicClient,
		resourceDiscovery: resourceDiscovery,
	}, nil
}

type LuaVM struct {
	l *lua.LState
}

func NewLuaVM() *LuaVM {
	L := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
	})
	registerKubeClientType(L)
	return &LuaVM{l: L}
}

func (l *LuaVM) Load(filename string) error {
	return l.l.DoFile(filename)
}

func (l *LuaVM) Close() {
	l.l.Close()
}
