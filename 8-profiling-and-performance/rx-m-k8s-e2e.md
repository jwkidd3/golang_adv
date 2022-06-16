# Go

# Docker
```
curl -L get.docker.com | sh
```

# kubectl
```
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
```

# Kind
```
go get sigs.k8s.io/kind@v0.11.1
export PATH=$(pwd)/go/bin:$PATH
kind version
```

```
~$ kind create cluster --name c1
Creating cluster "c1" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼 
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-c1"
You can now use your cluster with:

kubectl cluster-info --context kind-c1

Have a nice day! 👋
```

```
~$ kubectl version --client
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.2", GitCommit:"092fbfbf53427de67cac1e9fa54aaa09a28371d7", GitTreeState:"clean", BuildDate:"2021-06-16T12:59:11Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"linux/amd64"}
ubuntu@ip-172-31-42-126:~$ kubectl cluster-info --context kind-c1
Kubernetes control plane is running at https://127.0.0.1:36047
CoreDNS is running at https://127.0.0.1:36047/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
~$
```

# kubetest2
```
go get sigs.k8s.io/kubetest2/...@latest
```

# deployer (kind)
```
go get sigs.k8s.io/kubetest2/kubetest2-kind@latest
```

# testers (exec, ginkgo)
```
go get sigs.k8s.io/kubetest2/kubetest2-tester-exec@latest
go get sigs.k8s.io/kubetest2/kubetest2-tester-ginkgo@latest
```