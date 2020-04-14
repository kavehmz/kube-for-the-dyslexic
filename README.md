# kube-for-the-dyslexic
This is a series of demos to show how Kuberntes works.

What you need is an installation of `docker`, `kubectl` and a tool named `kind` which is used by Kubernetes team itself for testing.

Docker: https://docs.docker.com/get-docker/
Kubeectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/
Kind: https://kind.sigs.k8s.io/docs/user/quick-start/#installation


### Verify
Just to make sure you are all set lets try the following commands. You still dont need to understand anything.
Just relax and try.

Make sure your docker us working and you have a recent enough one.
```
$ docker version

# Make sure the following works fine.
# Here you want to run an image named hello-world in a container. This is a simple app and just prints some stuff.
# You docker will first check if image exists locally and if not docker tried to find it in ALL docker registry servers it knows.
$ docker run hello-world
docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
1b930d010525: Pull complete
Digest: sha256:f9dfddf63636d84ef479d645ab5885156ae030f611a56f3a7ac7f2fdd86d7e4e
Status: Downloaded newer image for hello-world:latest
Hello from Docker!
```

Verify your `kubectl`.
```
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"14+", GitVersion:"v1.14.10-dispatcher", GitCommit:"f5757a1dee5a89cc5e29cd7159076648bf21a02b", GitTreeState:"clean", BuildDate:"2020-02-06T03:31:35Z", GoVersion:"go1.12.12b4", Compiler:"gc", Platform:"darwin/amd64"}
The connection to the server localhost:8080 was refused - did you specify the right host or port?
```

Verify your `kind` setup.
```
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.17.0) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

# The kubectl cluster-info prints information about the control plane and add-ons. You will see later what they mean.
# For now we have the api server (as main component) and only KubeDNS (as addon)
$ kubectl cluster-info --context kind-kind
Kubernetes master is running at https://127.0.0.1:32772
KubeDNS is running at https://127.0.0.1:32772/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

# Do you remember we ran hello-world in docker? Kuberentes also uses containers to run stuff.
# Lets just run the same image in Kubernetes too. 
# Don't worry about details at all. Just try and get comfortable with the concepts.

# create a deployment that uses our image. Deployments set details like which image, how many replica, max memory, ... 
$ kubectl --context kind-kind create  deployment hello --image=hello-world
deployment.apps/hello created

#  But deployments are just setting the details. Our containers run inside Pods!
$ kubectl --context kind-kind get deployment -o wide
NAME    READY   UP-TO-DATE   AVAILABLE   AGE    CONTAINERS    IMAGES        SELECTOR
hello   0/1     1            0           3m4s   hello-world   hello-world   app=hello

# Pods are the resources in Kubertenes that actually show a group of contaierns running in one node (yes you can run multiple images in a single pod)
$ kubectl --context kind-kind get pod -owide
NAME                     READY   STATUS             RESTARTS   AGE     IP           NODE                 NOMINATED NODE   READINESS GATES
hello-67d96bb797-9rffr   0/1     CrashLoopBackOff   5          4m10s   10.244.0.6   kind-control-plane   <none>           <none>

# Ignore the CrashLoopBackOff and now look at the pods log
# You pod name will be different. 
# Mine is `hello-67d96bb797-9rffr`.
$ kubectl --context kind-kind logs -f hello-67d96bb797-9rffr
Hello from Docker!

# Lets push it a bit further. Still you just need to try and get comfortable with concepts and tools.
# this time lets describe and see what we find
$ kubectl --context kind-kind describe deployment hello
Name:                   hello
Namespace:              default
CreationTimestamp:      Tue, 14 Apr 2020 09:28:40 +0200
Labels:                 app=hello
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=hello
Replicas:               1 desired | 1 updated | 1 total | 0 available | 1 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=hello
  Containers:
   hello-world:
    Image:        hello-world
    Port:         <none>
    Host Port:    <none>
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      False   MinimumReplicasUnavailable
  Progressing    False   ProgressDeadlineExceeded
OldReplicaSets:  <none>
NewReplicaSet:   hello-67d96bb797 (1/1 replicas created)
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  15m   deployment-controller  Scaled up replica set hello-67d96bb797 to 1

# Describe says that deployment created a Replicaset named `hello-67d96bb797`.
# lets check it!
kubectl --context kind-kind get replicaset
NAME               DESIRED   CURRENT   READY   AGE
hello-67d96bb797   1         1         0       14m

# and describe it to get more info.
$ kubectl --context kind-kind describe replicaset
Name:           hello-67d96bb797
Namespace:      default
Selector:       app=hello,pod-template-hash=67d96bb797
Labels:         app=hello
                pod-template-hash=67d96bb797
Annotations:    deployment.kubernetes.io/desired-replicas: 1
                deployment.kubernetes.io/max-replicas: 2
                deployment.kubernetes.io/revision: 1
Controlled By:  Deployment/hello
Replicas:       1 current / 1 desired
Pods Status:    1 Running / 0 Waiting / 0 Succeeded / 0 Failed
Pod Template:
  Labels:  app=hello
           pod-template-hash=67d96bb797
  Containers:
   hello-world:
    Image:        hello-world
    Port:         <none>
    Host Port:    <none>
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Events:
  Type    Reason            Age   From                   Message
  ----    ------            ----  ----                   -------
  Normal  SuccessfulCreate  14m   replicaset-controller  Created pod: hello-67d96bb797-9rffr

# It might show now that we created a deployment `hello`, deployment created a repicaset `hello-67d96bb797` and replicaset created a pod `hello-67d96bb797-9rffr`
# later you will see one part of kuberetnes get the event for pod creation and pull the image and started the containers as it was indicated in pod details!.
# so our application ran in a node named `kind-control-plane`
$ kubectl --context kind-kind get pod
NAME                     READY   STATUS             RESTARTS   AGE
hello-67d96bb797-9rffr   0/1     CrashLoopBackOff   9          24m
```

# Lets check the setup now. You can see what containers are running using the following command.
$ docker ps
docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED              STATUS              PORTS                       NAMES
95623bda34ec        kindest/node:v1.17.0   "/usr/local/bin/entr…"   About a minute ago   Up About a minute   127.0.0.1:32772->6443/tcp   kind-control-plane




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
Read more: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/
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
Read more: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/
In Kubernetes, a controller is a control loop that watches the shared state of the cluster through the apiserver and makes changes attempting to move the current state towards the desired state, for example replication controller, endpoints controller, namespace controller, and serviceaccounts controller.

Document says things like:
* Node Controller: Responsible for noticing and responding when nodes go down.
* Replication Controller: Responsible for maintaining the correct number of pods for every replication controller object in the system.
* Endpoints Controller: Populates the Endpoints object (that is, joins Services & Pods).
* Service Account & Token Controllers: Create default accounts and API access tokens for new namespaces.

Lets try the `Node Controller` part and see how `kube-controller-manager` reacts when nodes go down.

```
# first, lets check our nodes and see there is only one
$ kubectl --context kind-kind  get node
NAME                 STATUS   ROLES    AGE     VERSION
kind-control-plane   Ready    master   6m29s   v1.17.0

# now to get comfortable with kind lets delete our cluster and create a new one with two nodes
$ kind delete cluster
Deleting cluster "kind" ...

$ kind create cluster --config 00_configs/kind_two_nodes.yaml
$ kubectl --context kind-kind  get node
NAME                 STATUS   ROLES    AGE   VERSION
kind-control-plane   Ready    master   88s   v1.17.0
kind-worker          Ready    <none>   50s   v1.17.0

# now lets delete the `kind-worker` node.
kubectl --context kind-kind  delete node kind-worker
node "kind-worker" deleted

# check the log for `kube-controller-manager`
$ kubectl --context kind-kind  -n kube-system logs --tail 6 kube-controller-manager-kind-control-plane
I0413 22:07:38.926126       1 event.go:281] Event(v1.ObjectReference{Kind:"Node", Namespace:"", Name:"kind-worker", UID:"e2dd5c77-7779-49a8-af07-c7ac781a21ce", APIVersion:"", ResourceVersion:"", FieldPath:""}): type: 'Normal' reason: 'RegisteredNode' Node kind-worker event: Registered Node kind-worker in Controller
I0413 22:09:48.965638       1 event.go:281] Event(v1.ObjectReference{Kind:"Node", Namespace:"", Name:"kind-worker", UID:"e2dd5c77-7779-49a8-af07-c7ac781a21ce", APIVersion:"", ResourceVersion:"", FieldPath:""}): type: 'Normal' reason: 'RemovingNode' Node kind-worker event: Removing Node kind-worker from Controller
I0413 22:10:38.881643       1 gc_controller.go:77] PodGC is force deleting Pod: kube-system/kube-proxy-jq7qk
I0413 22:10:38.887413       1 gc_controller.go:188] Forced deletion of orphaned Pod kube-system/kube-proxy-jq7qk succeeded
I0413 22:10:38.887452       1 gc_controller.go:77] PodGC is force deleting Pod: kube-system/kindnet-4mt6r
I0413 22:10:38.896510       1 gc_controller.go:188] Forced deletion of orphaned Pod kube-system/kindnet-4mt6r succeeded
```

### Control Plane Components: cloud-controller-manager
Read more: https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/
`cloud-controller-manager` runs controllers that interact with the underlying cloud providers.
These can overlap with `kube-controller-manager`.

## Node Components
Node components run on every node, maintaining running pods and providing the Kubernetes runtime environment.

#### dns
kubectl --context kind-kind -n kube-system port-forward coredns-6955765f44-r7drp 32053:53
dig @127.0.0.1 -p 32053 kube-dns.kube-system.svc.cluster.local +tcp
