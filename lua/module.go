package lua

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaKubeModuleName = "kube"

var kubeMethods = map[string]lua.LGFunction{
	"version":   kubeVersion,
	"resources": kubeResources,
}

func kubeVersion(L *lua.LState) int {
	k := checkKubeClient(L)
	v, err := k.clientset.Discovery().ServerVersion()
	if err != nil {
		L.RaiseError("failed to get kube version: %v", err)
		return 0
	}
	L.Push(lua.LString(v.String()))
	return 1
}

func kubeResources(L *lua.LState) int {
	k := checkKubeClient(L)
	resources := k.resourceDiscovery.List()
	resourcesTable := L.NewTable()
	for _, resource := range resources {
		gvrTable := L.NewTable()
		gvrTable.RawSetString("group", lua.LString(resource.Group))
		gvrTable.RawSetString("version", lua.LString(resource.Version))
		gvrTable.RawSetString("resource", lua.LString(resource.Resource))
		resourcesTable.Append(gvrTable)
	}
	L.Push(resourcesTable)
	return 1
}
