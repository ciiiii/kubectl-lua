package lua

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaKubeModuleName = "kube"

var kubeMethods = map[string]lua.LGFunction{
	"version": kubeVersion,
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
