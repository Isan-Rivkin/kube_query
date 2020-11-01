# kube_query

Use kubectl but on all of the available k8s clusters available in the kubeconfig file. 
Currently will query only AWS EKS clusters.

### Install 

```bash
git clone
make release
```
### Uninstall

```bash 
make clean
```

### Usage

It's a `kubectl` wrapper that querys all contexts and finished at the current one. 
So just use whatever kubectl command with that. 

```bash
kq <kubectl query>
```



### example

```bash
kq get pods --namespace kube-system
```