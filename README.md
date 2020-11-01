# kube_query

Use kubectl but on all of the available k8s clusters available in the kubeconfig file. 

### Install 

```bash
git clone 
go build main.go 
cp main.go <kq>
# on Macos to use kq globally 
mv kq /usr/local/bin/
```

### Usage

It's a `kubectl` wrapper that querys all contexts and finished at the current one. 
So just use whatever kubectl command with that. 

`kq <kubectl query>`

```bash 
```
### example

```bash
kq get pods -n kube-system
```