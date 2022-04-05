package vm

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/ciiiii/kubectl-lua/api"
)

type KubeClient struct {
	clientset         *kubernetes.Clientset
	dynamic           dynamic.Interface
	resourceDiscovery *api.ResourceDiscovery
}

func kubeClientFromConfig(config *rest.Config) (*KubeClient, error) {
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

func NewLuaVM(config *rest.Config) *LuaVM {
	L := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
	})
	registerKubeClientType(L, config)
	return &LuaVM{l: L}
}

func (l *LuaVM) Load(filename string) error {
	return l.l.DoFile(filename)
}

func (l *LuaVM) REPL() error {
	// TODO: implement REPL
	fmt.Println("REPL is not implemented yet")
	return nil
}

func (l *LuaVM) Close() {
	l.l.Close()
}
