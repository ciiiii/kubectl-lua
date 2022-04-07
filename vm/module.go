package vm

import (
	"context"

	lua "github.com/yuin/gopher-lua"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const LuaKubeModuleName = "kube"

var kubeMethods = map[string]lua.LGFunction{
	"version":   kubeVersion,
	"resources": kubeResources,
	"resource":  kubeResource,
	"namespace": kubeNamespace,
	"get":       kubeGet,
	"list":      kubeList,
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

func kubeResource(L *lua.LState) int {
	if L.GetTop() >= 2 {
		kube := checkKubeClient(L)
		kind := L.CheckString(2)
		newKube := *kube
		newKube.kind = kind
		userData := L.NewUserData()
		userData.Value = &newKube
		L.SetMetatable(userData, L.GetTypeMetatable(LuaKubeModuleName))
		L.Push(userData)
		return 1
	}
	L.RaiseError("arguments error: %d", L.GetTop())
	return 0
}

func kubeNamespace(L *lua.LState) int {
	if L.GetTop() == 2 {
		kube := checkKubeClient(L)
		namespace := L.CheckString(2)
		newKube := *kube
		newKube.namespace = namespace
		newKube.scope = namespaceScoped
		userData := L.NewUserData()
		userData.Value = &newKube
		L.SetMetatable(userData, L.GetTypeMetatable(LuaKubeModuleName))
		L.Push(userData)
		return 1
	}
	L.RaiseError("arguments error: %d", L.GetTop())
	return 0
}

func kubeGet(L *lua.LState) int {
	if L.GetTop() == 2 {
		k := checkKubeClient(L)
		name := L.CheckString(2)
		client, err := k.resourceClient()
		if err != nil {
			L.RaiseError("failed to init resource client: %v", err)
			return 0
		}
		obj, err := client.Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			L.RaiseError("failed to get resource: %v", err)
			return 0
		}
		val := decodeValue(L, obj.Object)
		L.Push(val)
		return 1
	}
	L.RaiseError("arguments error: %d", L.GetTop())
	return 0
}

func kubeList(L *lua.LState) int {
	if L.GetTop() == 1 {
		k := checkKubeClient(L)
		client, err := k.resourceClient()
		if err != nil {
			L.RaiseError("failed to init resource client: %v", err)
			return 0
		}
		objList, err := client.List(context.Background(), metav1.ListOptions{})
		if err != nil {
			L.RaiseError("failed to list resource: %v", err)
			return 0
		}
		tableData := L.NewTable()
		for _, obj := range objList.Items {
			val := decodeValue(L, obj.Object)
			tableData.Append(val)
		}

		L.Push(tableData)
		return 1
	}
	L.RaiseError("arguments error: %d", L.GetTop())
	return 0
}
