package vm

import (
	"context"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const LuaKubeModuleName = "kube"

var kubeMethods = map[string]lua.LGFunction{
	"version":      kubeVersion,
	"resources":    kubeResources,
	"listResource": kubeListResource,
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

func kubeListResource(L *lua.LState) int {
	k := checkKubeClient(L)
	if L.GetTop() == 4 {
		group := L.CheckString(2)
		version := L.CheckString(3)
		resource := L.CheckString(4)
		fmt.Printf("group: %s, version: %s, resource: %s\n", group, version, resource)
		objList, err := k.dynamic.Resource(schema.GroupVersionResource{
			Group:    group,
			Version:  version,
			Resource: resource,
		}).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			L.RaiseError(err.Error())
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
