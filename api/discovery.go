package api

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

type ResourceDiscovery struct {
	lister []schema.GroupVersionResource
	mapper map[string][]int
}

func NewResourceDiscovery(client discovery.DiscoveryInterface) (*ResourceDiscovery, error) {
	resources, err := client.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	lister := make([]schema.GroupVersionResource, 0, len(resources))
	mapper := make(map[string][]int)

	for _, resourceList := range resources {
		for _, apiResource := range resourceList.APIResources {
			gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
			if err != nil {
				continue
			}
			lister = append(lister, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: apiResource.Name,
			})
			index := len(lister) - 1
			names := allNames(apiResource, gv)
			for _, name := range names {
				if val, exist := mapper[name]; exist {
					mapper[name] = append(val, index)
				} else {
					mapper[name] = []int{index}
				}
			}
		}
	}
	return &ResourceDiscovery{
		lister: lister,
		mapper: mapper,
	}, nil
}

func (rd *ResourceDiscovery) gvrByIndex(index int) (schema.GroupVersionResource, error) {
	if index < 0 || index >= len(rd.lister) {
		return schema.GroupVersionResource{}, fmt.Errorf("index out of range")
	}
	return rd.lister[index], nil
}

func (rd *ResourceDiscovery) Search(kind string) (schema.GroupVersionResource, error) {
	if indexes, exist := rd.mapper[strings.ToLower(kind)]; exist {
		if len(indexes) == 1 {
			return rd.gvrByIndex(indexes[0])
		}
		resources := make([]string, 0, len(indexes))
		for _, index := range indexes {
			gvr, err := rd.gvrByIndex(index)
			if err == nil {
				resources = append(resources, strings.Join([]string{gvr.Resource, gvr.Version, gvr.Group}, "."))
			}
		}
		return schema.GroupVersionResource{}, fmt.Errorf("found multiple resources %s for kind %q", strings.Join(resources, ", "), kind)
	} else {
		return schema.GroupVersionResource{}, fmt.Errorf("kind %s not found", kind)
	}
}

func allNames(apiResource metav1.APIResource, gv schema.GroupVersion) []string {
	var baseNames []string
	baseNames = append(baseNames, apiResource.Name)
	baseNames = append(baseNames, apiResource.ShortNames...)
	if apiResource.SingularName != "" {
		baseNames = append(baseNames, apiResource.SingularName)
	} else {
		baseNames = append(baseNames, strings.ToLower(apiResource.Kind))
	}

	var fullNames []string
	for _, baseName := range baseNames {
		fullNames = append(fullNames,
			baseName,
			strings.Join([]string{baseName, gv.Group}, "."),
			strings.Join([]string{baseName, gv.Version, gv.Group}, "."),
		)
	}
	return fullNames
}
