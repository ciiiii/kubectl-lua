package vm

import (
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

var (
	namespaceScopedMethods = []string{"get", "list"}
	clusterScopedMethods   = append([]string{"resource", "resource", "namespace", "version"}, namespaceScopedMethods...)
)

func registerKubeClientType(L *lua.LState, config *rest.Config) {
	metaTable := L.NewTypeMetatable(LuaKubeModuleName)
	L.SetGlobal(LuaKubeModuleName, metaTable)
	L.SetField(metaTable, "new", L.NewFunction(newKubeClient(config)))
	L.SetField(metaTable, "__index", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 2 {
			k := checkKubeClient(L)
			index := L.CheckString(2)
			switch k.scope {
			case clusterScoped:
				for _, method := range clusterScopedMethods {
					if method == index {
						L.Push(L.NewFunction(kubeMethods[index]))
						return 1
					}
				}
				L.RaiseError("method or variable %q undefined", index)
				return 0
			case namespaceScoped:
				for _, method := range namespaceScopedMethods {
					if method == index {
						L.Push(L.NewFunction(kubeMethods[index]))
						return 1
					}
				}
				L.RaiseError("method or variable %q undefined", index)
				return 0
			default:
				L.RaiseError("invalid scope %q", k.scope)
				return 0
			}
		}
		L.RaiseError("arguments error: %d", L.GetTop())
		return 0
	}))
}

func newKubeClient(config *rest.Config) func(*lua.LState) int {
	return func(L *lua.LState) int {
		kubeClient, err := kubeClientFromConfig(config)
		if err != nil {
			L.RaiseError("failed to init kube client: %v", err)
			return 0
		}
		kubeClient.namespace = corev1.NamespaceAll
		kubeClient.scope = clusterScoped
		userData := L.NewUserData()
		userData.Value = kubeClient
		L.SetMetatable(userData, L.GetTypeMetatable(LuaKubeModuleName))
		L.Push(userData)
		return 1
	}
}

func checkKubeClient(L *lua.LState) *KubeClient {
	userData := L.CheckUserData(1)
	if kubeClient, ok := userData.Value.(*KubeClient); ok {
		return kubeClient
	}
	L.RaiseError("failed to get kube client from Lua userdata")
	return nil
}
