# kubectl-lua

A kubectl plugin to query and operate Kubernetes resources with Lua.

## Installation
TODO: submit plugin to krew index

## Usage
- run a lua file
```bash
kubectl-lua run <lua_file>
```
- start a lua interpreter
```bash
kubectl-lua repl
```

## API Reference
### kube.new
- Syntax: `kube.new()->kube`
- Type: Constructor
- Description: Create a new `kube` object.
- Return: A new `kube` object.
### kube:version
- Syntax: `kube:version()->string`
- Type: Method
- Description: Get the version of Kubernetes Cluster.
- Return: The version of Kubernetes Cluster.
### kube:resources
- Syntax: `kube:resources()->table`
- Type: Method
- Description: Get the Kubernetes api-resource list.
- Return: A table containing all GVR of Kubernetes resources.
### kube:listResource
- Syntax: `kube:listResource(group, version, kind)->table`
- Type: Method
- Description: List resources of target GVR.
- Return: A table containing all resources of target GVR.
