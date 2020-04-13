# kube-for-the-dyslexic
This is a series of demos to show how Kuberntes works.

What you need is an installation of `docker`.

We will build a docker image that includes all the required tools.

For Kubernetes part we will use [kind](https://kind.sigs.k8s.io/). Read the intro there.

### Docker
You need to install docker. Take a look at the [README](01_prepare_your_tools/README.md) if you want.

When you have docker ready, you can use the following command to build the image that has all you need.

```bash
$ make build
```

Tip: `make` command reads the [Makefile](./Makefile) and runs the section you mentioned, `build` in our case`. That is it.

If you like to see what you just build do this:

```bash
$ make run
```

You should be able to start your containers see be places inside it in a shell prompt. Try a echo command.

```bash
docker build -t local-k8s-tools -f 01_prepare_your_tools/Dockerfile . > /dev/null
docker run -ti --rm --name local-k8s-tools \
		-e HOME=/home/user \
		-v $PWD/00_data/home:/home/user \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $PWD:/workspace \
		local-k8s-tools
root@c6bc506827ce:/workspace# echo "Hellow World"
Hellow World
```

### Echo And Relay Apps
Now lets create two apps.
One named `echo` which is a simple http server that echos what we send send to it.
The other one is `relay` which relays our request to `echo` to serve it.

For now we just need two create a two docker images that conains the binary code of these two apps.
These two are written in Go language but we don't need anything but Docker again. We will compile these apps inside docker.

You can see how in this [Dockerfile](./00_apps/Dockerfile).

To create your images run the following make commands

echo:
```bash
$ make build-echo
docker build -t local-echo-server:latest -f 00_apps/Dockerfile 00_apps/echo_server
Sending build context to Docker daemon  5.167kB
Step 1/7 : FROM golang:1 AS build
 ---> 315fc470b445
Step 2/7 : COPY . /workspace
 ---> 874a7f57db4b
Step 3/7 : WORKDIR /workspace
 ---> Running in d9447a57c0e0
Removing intermediate container d9447a57c0e0
 ---> b6a24a684734
Step 4/7 : RUN go build -o /usr/local/bin/my_service main.go
 ---> Running in c3a17f33e1c3
Removing intermediate container c3a17f33e1c3
 ---> 10ccb2e280ac
Step 5/7 : FROM debian:stable-slim
 ---> e7e5f8b110eb
Step 6/7 : COPY --from=build /usr/local/bin/my_service /usr/local/bin/my_service
 ---> 31f6d01a35dc
Step 7/7 : ENTRYPOINT [ "/usr/local/bin/my_service" ]
 ---> Running in 27b75c6508dc
Removing intermediate container 27b75c6508dc
 ---> 7167ae5a3fe4
Successfully built 7167ae5a3fe4
Successfully tagged local-echo-server:latest
```

and relay:
```bash
make build-relay
docker build -t local-relay-server:latest -f 00_apps/Dockerfile 00_apps/relay_server
Sending build context to Docker daemon  5.679kB
Step 1/7 : FROM golang:1 AS build
 ---> 315fc470b445
Step 2/7 : COPY . /workspace
 ---> c987057e0d7f
Step 3/7 : WORKDIR /workspace
 ---> Running in 1b4f5aa3ab2a
Removing intermediate container 1b4f5aa3ab2a
 ---> 68717c19f285
Step 4/7 : RUN go build -o /usr/local/bin/my_service main.go
 ---> Running in 60b752124952
Removing intermediate container 60b752124952
 ---> 4a52d22e3ce6
Step 5/7 : FROM debian:stable-slim
 ---> e7e5f8b110eb
Step 6/7 : COPY --from=build /usr/local/bin/my_service /usr/local/bin/my_service
 ---> 8a4ec2b58c2b
Step 7/7 : ENTRYPOINT [ "/usr/local/bin/my_service" ]
 ---> Running in 37812c90f84b
Removing intermediate container 37812c90f84b
 ---> 854981eb023c
Successfully built 854981eb023c
Successfully tagged local-relay-server:latest
```

### Testing echo using Docker
Lets see how we can test `echo`.

We will run echo as a docker container so it runs isolated from out local environment.
Echo is a server that listens for http requests on port `8080`.
Notice we need to __expose__ that 8080 port from inside the contaier to outside.
For future reference remember that Kubernetes needs to do the same! We will try to see how in future steps.

Lets run `echo` server and test it.

```bash
$ make run-echo
docker run -ti -p 8081:8080 local-echo-server
fb7a32ccc83faf6ee3058cc12e957a9a03ab439e3d999282273013584459ee53
```

That large string you see above is our `CONTAINER ID`.

You can check our running container and some details using `docker ps` command

```
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED              STATUS              PORTS                    NAMES
fb7a32ccc83f        local-echo-server   "/usr/local/bin/my_s…"   About a minute ago   Up About a minute   0.0.0.0:8081->8080/tcp   condescending_nightingale
```

Pay attension to how `PORTS` sections indicates that we have a mapping from 0.0.0.0:8001 in local environemnt to 8080 port in container (on TCP).
We set an option `-p 8081:8080` for this in our `docker run`.

We can try and send a request. Our app inside container is listening to port 8080 but docker is exposing port 8081:

```
$ curl http://localhost:8081/echo?message=hello
Echo host: 8c68245116b9
Message: hello
```

Lets stop our container and move on:

```
make stop-echo
docker stop local-echo-server
local-echo-server
```

If you like start our container and map different ports and see what happens.

### Kubernetes Using kind
Now we need to create our Kubernetes cluster to run our two apps in there.

As said we will use [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/).
kind is a tool for running local Kubernetes clusters using Docker container “nodes”.

Verify if you have the kind command ready,

```bash
$ kind version
kind v0.7.0 go1.13.6 darwin/amd64
```

### Create your cluster
Creating a kuberentes clsuter using kind is simple

```bash
$ make create-cluster
kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.17.0) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

Now we can inspect our setup.

```bash
$ docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED              STATUS              PORTS                       NAMES
ea58ca5b4aaa        kindest/node:v1.17.0   "/usr/local/bin/entr…"   About a minute ago   Up About a minute   127.0.0.1:32772->6443/tcp   kind-control-plane
```

and lets see if our kubectl command works now.

```bash
$ kubectl --context kind-kind get all --all-namespaces
NAMESPACE            NAME                                             READY   STATUS    RESTARTS   AGE
kube-system          pod/coredns-6955765f44-r7drp                     1/1     Running   0          3m11s
kube-system          pod/coredns-6955765f44-sz8sz                     1/1     Running   0          3m11s
kube-system          pod/etcd-kind-control-plane                      1/1     Running   0          3m25s
kube-system          pod/kindnet-trhvj                                1/1     Running   0          3m11s
kube-system          pod/kube-apiserver-kind-control-plane            1/1     Running   0          3m25s
kube-system          pod/kube-controller-manager-kind-control-plane   1/1     Running   0          3m25s
kube-system          pod/kube-proxy-wxstw                             1/1     Running   0          3m11s
kube-system          pod/kube-scheduler-kind-control-plane            1/1     Running   0          3m25s
local-path-storage   pod/local-path-provisioner-7745554f7f-m9t4z      1/1     Running   0          3m11s

NAMESPACE     NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
default       service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP                  3m28s
kube-system   service/kube-dns     ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   3m26s

NAMESPACE     NAME                        DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                 AGE
kube-system   daemonset.apps/kindnet      1         1         1       1            1           <none>                        3m25s
kube-system   daemonset.apps/kube-proxy   1         1         1       1            1           beta.kubernetes.io/os=linux   3m26s

NAMESPACE            NAME                                     READY   UP-TO-DATE   AVAILABLE   AGE
kube-system          deployment.apps/coredns                  2/2     2            2           3m26s
local-path-storage   deployment.apps/local-path-provisioner   1/1     1            1           3m24s

NAMESPACE            NAME                                                DESIRED   CURRENT   READY   AGE
kube-system          replicaset.apps/coredns-6955765f44                  2         2         2       3m11s
local-path-storage   replicaset.apps/local-path-provisioner-7745554f7f   1         1         1       3m11s
```

As you see that single container created a lot of stuff. Lets explain it a bit.

### How Kubernetes Works
Read more here: https://kubernetes.io/docs/concepts/overview/components/

But in general we have three parts:
* `Control Plane Components`, The Control Plane’s components make global decisions about the cluster (for example, scheduling)
* `Node Components`, Node components run on every node.
* `Addons` 


#### Control Plane Components: kube-apiserver
Read more here: https://kubernetes.io/docs/concepts/overview/components/#kube-apiserver

The API server is a component of the Kubernetes control plane that exposes the Kubernetes API.

Lets run `kubectl` command again:

```
$ kubectl --context kind-kind --kubeconfig ~/.kube/config get namespaces
NAME                 STATUS   AGE
default              Active   68m
kube-node-lease      Active   68m
kube-public          Active   68m
kube-system          Active   68m
local-path-storage   Active   67m
```

Almost all kubectl does is getting the config from `--kubeconfig ~/.kube/config` and connecting to the cluster we want `--context kind-kind` which its address is defined in config file, and sending an API command to `kube-apiserver`! That is it.

Btw if you want to know what types of api calls are available in your cluster try this:

```bash
$ kubectl --context kind-kind api-resources -o wide
NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND                             VERBS
configmaps                        cm                                          true         ConfigMap                        [create delete deletecollection get list patch update watch]
namespaces                        ns                                          false        Namespace                        [create delete get list patch update watch]
pods                              po                                          true         Pod                              [create delete deletecollection get list patch update watch]
.....
```

It is useful if you check the `Kind` for each api, `VERBS` you can use and if it is `NAMESPACED` or not.

Later you will see how to extend Kuberentes API.

#### Control Plane Components: etcd
Read more: https://etcd.io/docs

This is an other component which us used as a key value store in as Kubernetes.

Lets try something here. (Notice all these components are in `kube-system`!)

```bash
$ kubectl --context kind-kind -n kube-system  get pod -owide|grep etcd
NAME                                         READY   STATUS    RESTARTS   AGE    IP           NODE                 NOMINATED NODE   READINESS GATES
etcd-kind-control-plane                      1/1     Running   0          5h7m   172.17.0.3   kind-control-plane   <none>           <none>
...
```

You can execute a command in a container using `kubectl exec PODNAME COMMAND`. Here we can try the following

```
$ kubectl --context kind-kind -n kube-system exec ti etcd-kind-control-plane /bin/sh
# Now you are in a shell inside etcd container.
# Lets see use the etcd key/value ability ourself.
# put KEY VALUE
$ etcdctl \
--cert /etc/kubernetes/pki/etcd/peer.crt \
--key /etc/kubernetes/pki/etcd/peer.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--endpoints https://127.0.0.1:2379 put some_key some_value

# Now lets get the value!
# get KEY
$ etcdctl \
--cert /etc/kubernetes/pki/etcd/peer.crt \
--key /etc/kubernetes/pki/etcd/peer.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--endpoints https://127.0.0.1:2379 get some_key
some_key
some_value

# just for fun lets also get what Kubernetes itself saved about this pod!
# (Pods are the smallest deployable units of computing that can be created and managed in Kubernetes)
# structure of key is intuitive
$ etcdctl \
--cert /etc/kubernetes/pki/etcd/peer.crt \
--key /etc/kubernetes/pki/etcd/peer.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--endpoints https://127.0.0.1:2379 get /registry/pods/kube-system/etcd-kind-control-plane -w json
{"header":{"cluster_id":9676036482053611986,"member_id":12858828581462913056,"revision":31446,"raft_term":2},"kvs":[{"key":"L3JlZ2lzdHJ5L3BvZHMva3ViZS1zeXN0ZW0vZXRjZC1raW5kLWNvbnRyb2wtcGxhbmU=","create_revision":249,"mod_revision":284,"version":2,"value":"......=="}],"count":1}
```

We wont need more than that I think. If you ever needed to create and maintain you own Kubernetes cluster just make yourself familiar with this tool and its administration.
Otherwise it is unlikely you will need to interact with etcd directly.


### Control Plane Components: kube-scheduler
Read more: https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/
Control plane component that watches for newly created Pods with no assigned node, and selects a node for them to run on.

Lets explore what Pod means first. And also lets try concept of namespace.

First, in kuberentes we are using containers to run our services. In our case we have a docker image named  that at the moment only exists in the local docker.  

Before you try these command check content of the yaml files first. You might find the comments there useful.
[Namepsace](./00_apps/echo_server/k8s/my_space.yaml)
[Pod](./00_apps/echo_server/k8s/pod.yaml)
```bash
# Lets create a namespace named `my-space`. check content of the yaml file first
$ kubectl --context kind-kind apply -f ./00_apps/echo_server/k8s/my_space.yaml
namespace/my-space created

# Lets create a pod now.
$ kubectl --context kind-kind apply -f ./00_apps/echo_server/k8s/pod.yaml
pod/echo created
```

Now lets see what happened.

```
$ kubectl --context kind-kind --namespace my-space get pod
NAME   READY   STATUS             RESTARTS   AGE
echo   0/1     ImagePullBackOff   0          1m

# error! let get some more info
$ kubectl --context kind-kind --namespace my-space describe pod echo
......
  Warning  Failed     12m (x41 over 3h41m)    kubelet, kind-control-plane  Error: ImagePullBackOff
  Normal   BackOff    2m48s (x85 over 3h40m)  kubelet, kind-control-plane  Back-off pulling image "local-echo-server:latest"
```

First, in kuberentes we are using containers to run our services. In our case we have a docker image named `local-echo-server:latest` that at the moment only exists in the local docker. That image is required by the container named `echo-container` (look at the yaml file to find both names).
Second image we use need to be accessible to the cluster! In nonral situations it means a docker registry like [docker hub](https://hub.docker.com/) or a private one that needs login.

In our case (which is special to kind) we have a command
```
kind load docker-image local-echo-server:latest
Image: "local-echo-server:latest" with ID "sha256:7c6d9d3501eacf69cfe4a1bc9abdd977c366702278ef399f643421f10ee2bb54" not present on node "kind-control-plane"
```

And now lets check  our state again:
```
$ kubectl --context kind-kind --namespace my-space get pod -o wide
NAME   READY   STATUS    RESTARTS   AGE    IP           NODE                 NOMINATED NODE   READINESS GATES
echo   1/1     Running   0          5m4s   10.244.1.7   kind-control-plane   <none>           <none>
```

This time we have one pod `RUNNING` and `READY`.

Notice between the time we were getting `ImagePullBackOff` error and next that we had a running pod we just made the image accessible.
`kube-scheduler` first noticed we have a new `Pod`. It schedules our pod in a node named `kind-control-plane`. 
Later we will see the another part of our setup named `kubelet` does a bit more and actually started the container on that node.

### Control Plane Components: kube-controller-manager

#### dns
kubectl --context kind-kind -n kube-system port-forward coredns-6955765f44-r7drp 32053:53
dig @127.0.0.1 -p 32053 kube-dns.kube-system.svc.cluster.local +tcp
