package lua

import lua "github.com/yuin/gopher-lua"

func registerKubeClientType(L *lua.LState) {
	metaTable := L.NewTypeMetatable(LuaKubeModuleName)
	L.SetGlobal(LuaKubeModuleName, metaTable)
	L.SetField(metaTable, "new", L.NewFunction(newKubeClient))
	L.SetField(metaTable, "__index", L.SetFuncs(L.NewTable(), kubeMethods))
}

func newKubeClient(L *lua.LState) int {
	kubeClient, err := kubeClientFromConfig()
	if err != nil {
		L.RaiseError("failed to init kube client: %v", err)
		return 0
	}
	userData := L.NewUserData()
	userData.Value = kubeClient
	L.SetMetatable(userData, L.GetTypeMetatable(LuaKubeModuleName))
	L.Push(userData)
	return 1
}

func checkKubeClient(L *lua.LState) *KubeClient {
	userData := L.CheckUserData(1)
	if kubeClient, ok := userData.Value.(*KubeClient); ok {
		return kubeClient
	}
	L.RaiseError("failed to get kube client from Lua userdata")
	return nil
}
