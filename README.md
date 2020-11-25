# kube_query

Use kubectl but on all of the available k8s clusters available in the kubeconfig file. 
Currently will query only AWS EKS clusters.

## Install 

### Brew 

MacOS (And ubuntu supported) installation via brew:

```bash
brew tap isan-rivkin/toolbox
brew install kq
```

### from source 

```bash
git clone
# default location on for Macos /usr/local/bin
make release
# or use make release BIN_DIR='/path/to/kq/binary' for custom kq binary location
```

### Target OS: 

```bash

make build-linux

make build-osx

make build-windows
```

### Usage

It's a `kubectl` wrapper that querys all contexts (e.g all of your clusters).
Just use whatever kubectl command with that.  

```bash
kq <kubectl query>
```



### example

```bash
kq get pods --namespace kube-system
kq get ns
kq describe pod --namespace default
```
