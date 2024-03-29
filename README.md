# kube-for-the-dyslexic
This readme is a step-by-step explanation of how Kubernetes works.

I tried to write it in a way that doesn't need you to memorize anything.

You don't need to concentrate. Read through and try the examples.

I will repeat the concepts again and again through this readme.

You need to follow the steps and absorb as much as you can.

You need an installation of `docker, `kubectl`, and a tool named `kind` which can create a Kubernetes cluster in docker.

Docker: https://docs.docker.com/get-docker/

kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/

Kind: https://kind.sigs.k8s.io/docs/user/quick-start/#installation


## Verify Your Tools

Ensure your docker is working and you have a recent enough version (>18).

```bash
$ docker version
Client: Docker Engine - Community
Version:           19.03.8
```

Make sure the following works fine.
Here you want to run an image named `hello-world` in a container. `hello-world` is a simple app that prints some stuff.
Your docker will first check if the image exists locally; if not, docker tries to find it in the registry servers it knows.

```
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
1b930d010525: Pull complete
Digest: sha256:f9dfddf63636d84ef479d645ab5885156ae030f611a56f3a7ac7f2fdd86d7e4e
Status: Downloaded newer image for hello-world:latest
Hello from Docker!
```

If you like to get more info about your docker setup use `docker info`.

```
$ docker info
Client:
 Debug Mode: false

Server:
 Containers: 9
  Running: 1
  Paused: 0
  Stopped: 8
 Images: 22
 Server Version: 19.03.8
 Storage Driver: overlay2
  Backing Filesystem: <unknown>
  Supports d_type: true
  Native Overlay Diff: true
 Logging Driver: json-file
 Cgroup Driver: cgroupfs
 Plugins:
  Volume: local
  Network: bridge host ipvlan macvlan null overlay
  Log: awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog
...
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
```

Lets check the setup now. You can see what containers are running using `docker`.
```
$ docker ps
docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED              STATUS              PORTS                       NAMES
95623bda34ec        kindest/node:v1.17.0   "/usr/local/bin/entr…"   About a minute ago   Up About a minute   127.0.0.1:32772->6443/tcp   kind-control-plane
```

kind started one container named kind-control-plane. This will act as a Kubernetes node for kind.
In Kubernetes we create several nodes (in this case only one) and start main Kubernetes components on them. In a more standard setup you will have more nodes. Some nodes will only run Kubernetes control plane components and some nodes for running other work loads.

Lets delete the cluster now. `kind` is a very light tool so don't be afraid of recreating your cluster to have a new one.

```
$ kind delete cluster
Deleting cluster "kind" ...
```

This time we create a new cluster but by giving `kind` a [config file](./00_configs/kind_two_nodes.yaml). Look inside the config file and see how this time we are instructing `kind` to create two nodes.

```
$ kind create cluster --config 00_configs/kind_two_nodes.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.17.0) 🖼
 ✓ Preparing nodes 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜

# lets check the docker containers again
$ docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                       NAMES
944688070a50        kindest/node:v1.17.0   "/usr/local/bin/entr…"   2 minutes ago       Up 2 minutes        127.0.0.1:32774->6443/tcp   kind-control-plane
d194a2e44c85        kindest/node:v1.17.0   "/usr/local/bin/entr…"   2 minutes ago       Up 2 minutes                                    kind-worker

# Lets check kubectl report about nodes too!
$ kubectl --context kind-kind  get node
NAME                 STATUS   ROLES    AGE     VERSION
kind-control-plane   Ready    master   3m20s   v1.17.0
kind-worker          Ready    <none>   2m48s   v1.17.0
```

The `kubectl cluster-info` prints information about the control-plane and add-ons. You will see later what they mean. For now we have the api server (as main component) and only KubeDNS (as addon)
```
$ kubectl cluster-info --context kind-kind
Kubernetes master is running at https://127.0.0.1:32772
KubeDNS is running at https://127.0.0.1:32772/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```

If all above was OK we can try running `hello-world` but this time in Kubernetes. Kubernetes also uses containers to run stuff. Don't worry about the details. We will go through all these details. Just try get the feel of command and concepts.

```
# Create a deployment that uses our image. Deployments set details like which image, how many replica, max memory, ... 
$ kubectl --context kind-kind create  deployment hello --image=hello-world
deployment.apps/hello created

# Deployments set the details.
$ kubectl --context kind-kind get deployment -o wide
NAME    READY   UP-TO-DATE   AVAILABLE   AGE    CONTAINERS    IMAGES        SELECTOR
hello   0/1     1            0           3m4s   hello-world   hello-world   app=hello

# But our containers run inside Pods!
# Pods are the resources in Kubernetes that present a group of containers running in one node (yes you can run multiple images in a single pod)
$ kubectl --context kind-kind get pod -o wide
NAME                     READY   STATUS             RESTARTS   AGE     IP           NODE                 NOMINATED NODE   READINESS GATES
hello-67d96bb797-9rffr   0/1     CrashLoopBackOff   5          4m10s   10.244.0.6   kind-control-plane   <none>           <none>

# Ignore the CrashLoopBackOff and look at the pods log
# Your pod name will be different. Mine is `hello-67d96bb797-9rffr`.
$ kubectl --context kind-kind logs -f hello-67d96bb797-9rffr
Hello from Docker!

# Lets push it a bit further.
# This time lets use `describe` and see what we can find
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

# And describe it to get more info.
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

# Last describe shows our repicaset `hello-67d96bb797` created a pod named `hello-67d96bb797-9rffr`.
# So Deployment created ReplicaSet which in turn created our Pod.
# Later you will see the next part of Kuberetnes flow which a service named `kubelet` uses the Pod definition and pulls our image and starts our container as it was indicated in Pod details.
$ kubectl --context kind-kind get pod
NAME                     READY   STATUS             RESTARTS   AGE
hello-67d96bb797-9rffr   0/1     CrashLoopBackOff   9          24m

# Delete the kind cluster we just created
# Create/Delete are light actions for kind. We can always start with a new cluster if we like
$ kind delete cluster
Deleting cluster "kind" ...
```

### Creating a Docker Image
Now lets create a docker image and run it in a container to have a feel for how it works.

We will create an image named `local-echo` which includes a simple http server that echos what we send send to it.

In docker images are defined by Dockerfile which is a list of steps, starting `FROM` a base image and then step by step `COPY` or `RUN` commands in that base image to change it. For example starting `FROM` a `golang:1` image naming it `build`, then `COPY`ing some Go files and compiling them to have a binary file. Then switching to a new image `debian:stable-slim` to start clean and just copying the binary file from previous steps. the result will be our image which we will use to run our app.

You can see the details in this [Dockerfile](./00_apps/Dockerfile). It is easy to understand.
Read more: https://docs.docker.com/engine/reference/builder/

To create your images run the following command:

echo:
```bash
# local-echo-server:latest is full name of the image (repository:tag)
$ docker build -t local-echo-server:latest -f 00_apps/Dockerfile 00_apps/echo_server
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

# List the images that you have locally
$ docker images
docker images
REPOSITORY                            TAG                 IMAGE ID            CREATED              SIZE
<none>                                <none>              52c9ea485687        About a minute ago   817MB
local-echo-server                     latest              46e818efa1f7        9 days ago           76.8MB
golang                                1                   315fc470b445        2 weeks ago          809MB
debian                                stable-slim         e7e5f8b110eb        3 weeks ago          69.2MB
kindest/node                          <none>              ec6ab22d89ef        3 months ago         1.23GB
hello-world                           latest              fce289e99eb9        15 months ago        1.84kB
```


### Running your container in docker
Lets see how we can test our `echo` app.

We will run echo as a docker container so it runs isolated from out local environment.
Echo is a server that listens for http requests on port `8080`.
Notice we need to __expose__ that 8080 port from inside the container to outside.
Notice that Kubernetes needs to do the same! We will see how in future steps.

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

Pay attention to how `PORTS` sections indicates that we have a mapping from 0.0.0.0:8001 in local environment to 8080 port in container (on TCP).
We set an option `-p 8081:8080` for this in our `docker run`.

We can try and send a request. Our app inside container is listening to port 8080 but docker is mapping port 8081 in our local environment to port 8080 inside the container.:

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

### Create your cluster
Creating a Kubernetes cluster using kind is simple

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

Let's see what kubectl shows.

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

As you see that single container created a lot of stuff. In next steps you will see what those are.

## How Kubernetes Works
Read more here: https://kubernetes.io/docs/concepts/overview/components/

But in general we have three parts:
* `Control Plane Components`, The Control Plane’s components make global decisions about the cluster (for example, scheduling)
* `Node Components`, Node components run on every node.
* `Addons` 


### Control Plane Components: kube-apiserver
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

Later you will see how to extend Kubernetes API.

### Control Plane Components: etcd
Read more: https://etcd.io/docs

This is an other component which us used as a key value store in as Kubernetes.

Lets try something here. (Notice all these components are in `kube-system`!)

```bash
$ kubectl --context kind-kind -n kube-system  get pod -o wide|grep etcd
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

First, in Kubernetes we are using containers to run our services. In our case we have a docker image named  that at the moment only exists in the local docker.

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

First, in Kubernetes we are using containers to run our services. In our case we have a docker image named `local-echo-server:latest` that at the moment only exists in the local docker. That image is required by the container named `echo-container` (look at the yaml file to find both names).
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

# now to get comfortable with `kind` lets delete our cluster and create a new one with two nodes
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

### Node Components: kubelet
An agent that runs on each node in the cluster. It makes sure that containers are running in a Pod.

Lets check all these:
```
# first lets delete any current Kubernetes cluster in `kind`
$ kind delete cluster
Deleting cluster "kind" ...

# then lets create a two node cluster
$ kind create cluster
$ kubectl --context kind-kind get pod --all-namespaces
NAMESPACE     NAME                                         READY   STATUS    RESTARTS   AGE
kube-system   etcd-kind-control-plane                      1/1     Running   0          5s
kube-system   kube-apiserver-kind-control-plane            1/1     Running   0          5s
kube-system   kube-controller-manager-kind-control-plane   1/1     Running   0          5s
kube-system   kube-scheduler-kind-control-plane            1/1     Running   0          5s
```

First notice there is no pod named kubelet. kubelets is the component that is responsible for running pods so unlike many other components kubelets themselves are not install as Pods.

So we will do an experiment. We will create a Pod, this time for running nginx, first while kubelet is not running and then we will start kubelet and see what happens.

```
kubectl --context kind-kind apply -f <(echo '{"apiVersion": "v1","kind": "Pod","metadata": {"name": "nginx"},"spec": {"containers": [{"name": "nginx","image": "nginx"}]}}')
pod/ngnix created

kubectl --context kind-kind describe pod nginx

docker exec -ti kind-control-plane /bin/bash -c 'systemctl status kubelet'
docker exec -ti kind-control-plane /bin/bash -c 'systemctl stop kubelet'

kubectl --context kind-kind get nodes
kubectl --context kind-kind describe node kind-control-plane
```


### Node Components: kube-proxy
kube-proxy is a network proxy that runs on each node in your cluster, implementing part of the Kubernetes Service concept.

kube-proxy maintains network rules on nodes. These network rules allow network communication to your Pods from network sessions inside or outside of your cluster.

kube-proxy uses the operating system packet filtering layer if there is one and it’s available. Otherwise, kube-proxy forwards the traffic itself

### Node Components: Container Runtime
The container runtime is the software that is responsible for running containers.

Kubernetes supports several container runtimes: Docker, containerd, CRI-O, and any implementation of the Kubernetes CRI (Container Runtime Interface).

Here we are doing it a bit nested in `kind`.

`kind create cluster` creates Kubernetes cluster by using docker to start containers and each containers represend one Node in our cluster.

```
$ docker ps 
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                       NAMES
7cf0a3b6e82b        kindest/node:v1.17.0   "/usr/local/bin/entr…"   16 hours ago        Up 16 hours         127.0.0.1:32776->6443/tcp   kind-control-plane
```

Then in those nodes, for example `kind-control-plane`, in the running debin based container, our image has `containerd` installed which runs our containers inside `kind-control-plane`.

Lets see what is running inside our node
```
$ docker  exec -ti kind-control-plane /bin/bash -c 'ps aux|grep container'
root     67673  1.7  1.5 2153040 61980 ?       Ssl  11:47   0:07 /usr/local/bin/containerd
```

And if you want to list the running containers inside `kind-control-plane` you can try
```
# inside the container there is a different tools installed for interaction with containerd. It is `crictl`. 
# https://github.com/kubernetes-sigs/cri-tools/blob/master/docs/crictl.md
$ docker  exec -ti kind-control-plane /bin/bash -c 'crictl ps '
CONTAINER           IMAGE               CREATED             STATE               NAME                      ATTEMPT             POD ID
b99d95197df05       5a8dfb2ca7312       7 minutes ago       Running             nginx                     1                   dc3e5938a2c6a
1fe5a9591be65       70f311871ae12       16 hours ago        Running             coredns                   0                   0ec4a02f58b0e
62be4c1062748       9d12f9848b99f       16 hours ago        Running             local-path-provisioner    0                   4629dc92cdd6d
fff092b0ccd01       70f311871ae12       16 hours ago        Running             coredns                   0                   575c12113852d
8a13f67aa95b3       2186a1a396deb       16 hours ago        Running             kindnet-cni               0                   8eb7a258a3ffe
e2bba59647f43       551eaeb500fda       16 hours ago        Running             kube-proxy                0                   ea35384f77c7d
6c7ff0a980651       303ce5db0e90d       16 hours ago        Running             etcd                      0                   98339fdb5233f
c7b024bf42b37       134ad2332e042       16 hours ago        Running             kube-apiserver            0                   3f35a132350a3
656999c9ec16e       09a204f38b41d       16 hours ago        Running             kube-scheduler            0                   a6953c193d5c4
9808ba4a683f0       7818d75a7d002       16 hours ago        Running             kube-controller-manager   0                   f56480b04bee2
```

## Addons
Addons use Kubernetes resources (DaemonSet, Deployment, etc) to implement cluster features. Because these are providing cluster-level features, namespaced resources for addons belong within the kube-system namespace.


#### Addons: DNS
Read more: https://kubernetes.io/docs/concepts/overview/components/#dns

While the other addons are not strictly required, all Kubernetes clusters should have cluster DNS, as many examples rely on it.

Cluster DNS is a DNS server, in addition to the other DNS server(s) in your environment, which serves DNS records for Kubernetes services.

Lets find our DNS server. 

```
# Here instead of just getting all pods and using a `grep` to filter for dns we are getting all pods which are labelled as `k8s-app=kube-dns`
# You can check definition of pods and their labels using `describe pod pod_name` command if you don't know the labels for a specific pod
$ kubectl --context kind-kind -n kube-system get pod -o wide -l k8s-app=kube-dns
NAME                       READY   STATUS    RESTARTS   AGE   IP           NODE                 NOMINATED NODE   READINESS GATES
coredns-6955765f44-7nqw2   1/1     Running   0          16h   10.244.0.2   kind-control-plane   <none>           <none>
coredns-6955765f44-b7szn   1/1     Running   0          16h   10.244.0.4   kind-control-plane   <none>           <none>
```

Now lets see why we say there is our dns servers. To check their functionality we will do a `port-forward` which opens a `tcp` (at tis moment no udp) connection from our local environemnt to the pod in Kubernetes cluster. It is a useful investigation tool.

```
# just notice your pod name will be different. Also later you will see we could do the same thing using DNS service.
$ kubectl --context kind-kind -n kube-system port-forward coredns-6955765f44-7nqw2 32053:53
Forwarding from 127.0.0.1:32053 -> 53
Forwarding from [::1]:32053 -> 53
```

Here we forward a local port `32053` to Pod port `53`.

Now open a new terminal and try the following:
```
# here we use a dig command to query our dns server, at local machine (address 127.0.0.1), at port 32053 (otherwise default port is 53), using tcp protocol (port-forward does not support udp yet).
# we are asking for IP of github.com
$ dig @127.0.0.1 -p 32053 +tcp github.com
;; ANSWER SECTION:
github.com.		22	IN	A	140.82.113.3

# You can try another dns server like CLoudflare's one just to compare.
# This time we dont need to ask for tcp protocol or a non standard port. so command is shorter
$ dig @1.1.1.1 github.com
;; ANSWER SECTION:
github.com.		19	IN	A	140.82.114.3
```

But that is not all our DSN server in the cluster does. It also resolve the address for services we have inside the cluster.
Actually that is what our DSN server is installed!

But first lets find some services. We will learm more about what services are later.

```
$kubectl --context kind-kind -n default get services
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   16h
```

We can find one service in `default` namespace. It is a service for our api-server.
Its IP inside cluster is `10.96.0.1`. But this IP could be different and can change!
Our applications inside Kubernetes might not be able to rely on static IP to find each other.
So there is a naming convention.

Read more: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/

For exmaple if we want to find IP of a servive (A record), the name is as `my-svc.my-namespace.svc.cluster-domain.example`

In our case it is ` `

Lets see if our dsn server which we port-forwarded can tell us what IP `kubernetes.default.svc.cluster.local` points to.

```
$ dig @127.0.0.1 -p 32053 kube-dns.kube-system.svc.cluster.local +tcp
;; ANSWER SECTION:
kube-dns.kube-system.svc.cluster.local.	30 IN A	10.96.0.10
```

Notice Cloudflare DSN server has no clue about dns records inside our cluster.

```
$ dig @1.1.1.1 kube-dns.kube-system.svc.cluster.local
# You wont get any IP here!
```

Tip: If you are wondering how we created `kubernetes.default.svc.cluster.local` from `my-svc.my-namespace.svc.cluster-domain.example` format, this how:
my-svc: name of the service you saw in `kubectl --context kind-kind -n default get services`. It was `kubernetes`.
my-namespace: the namespace that service exists there. In our command we were searching `default` namespace. So our service is there.
svc: this is just an string. copy it there.
cluster-domain.example: This one can be set for clustsers. By default it is `cluster.local` but it might be different in different setups. It is part of the dsn config! 

```
# Kubernetes applications offers a resource that we can to store configureation in plain text and use that config in our Pods.
# It is called ConfigMap. The ConfigMap for `coredns` can we check using:
$ kubectl --context kind-kind -n kube-system get configmaps coredns -o yaml
kubectl --context kind-kind -n kube-system get configmaps coredns -o yaml
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
kind: ConfigMap

# Notice many setting of our DSN server is in ConfigMap like the port it listens to, ttl and of course you can find `cluster.local` there.
# If you like you can edit this configMap (and in general many resource which support `edit` command):
$ kubectl --context kind-kind -n kube-system edit configmaps coredns
```

#### Addons: Web UI
Read more: https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

Dashboard is a general purpose, web-based UI for Kubernetes clusters. It allows users to manage and troubleshoot applications running in the cluster, as well as the cluster itself.

`kind` won't install the Web UI and I also never found much use for it, but lets install and try it.

```
$ kubectl --context kind-kind apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta8/aio/deploy/recommended.yaml
namespace/kubernetes-dashboard created
serviceaccount/kubernetes-dashboard created
service/kubernetes-dashboard created
secret/kubernetes-dashboard-certs created
secret/kubernetes-dashboard-csrf created
secret/kubernetes-dashboard-key-holder created
configmap/kubernetes-dashboard-settings created
role.rbac.authorization.k8s.io/kubernetes-dashboard created
clusterrole.rbac.authorization.k8s.io/kubernetes-dashboard created
rolebinding.rbac.authorization.k8s.io/kubernetes-dashboard created
clusterrolebinding.rbac.authorization.k8s.io/kubernetes-dashboard created
deployment.apps/kubernetes-dashboard created
service/dashboard-metrics-scraper created
deployment.apps/dashboard-metrics-scraper created
```

It might be useful if you take a look at the yaml file content. Note it will create a namespace nemd `kubernetes-dashboard` and installing all resources in that namespace.

To access it you need to use `kubectl proxy` command.

```
kkubectl --context kind-kind proxy --port=8001
```
Read more: https://kubernetes.io/docs/tasks/access-kubernetes-api/http-proxy-access-api/

`proxy` comamnd authenticates you using the details in your kubeconfig, `~/.kube/config`. This is required whenever you want to communicate with Kuberenetes API.

This gives you acceess to the API now. Before we get into Web UI lets try few calls ourselves!

```
# get all namespaces
curl -X GET 'http://localhost:8001/api/v1/namespaces
# get details of `default` namespace
curl -X GET 'http://localhost:8001/api/v1/namespaces/default
# get pods in `kube-system` namespace
curl -X GET 'http://localhost:8001/api/v1/namespaces/kube-system/pods'
```

Notice all kubectl does is reading kubeconfig, authenticating with Kubernetes API and then sending similar command we just send using curl!

Now lets check the Web UI. In your browser open the following URL:

http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy

Notice we are asked for `Sign in`. Despite the fact that `kubectl` authenticated us using `kubeconfig` credentials, `kubectl` only connected us to the Kubernetes API so we can open this page.

There is another concept that we will learn about it in details later but lets get a gimpse of it now, ServiceAccounts.
A service account provides an identity for processes that run in a Pod.

```
# Lets create a new and empty namespace names `space`
kubectl --context kind-kind create namespace space
namespace/space created

# We will create a new service account named `webui` in `space`
kubectl --context kind-kind -n space create serviceaccount webui

# Every service account we create adds a sercret that contains its details
$ kubectl --context kind-kind -n space get serviceaccount
NAME      SECRETS   AGE
default   1         78s
webui     1         45s
$ kubectl --context kind-kind -n space get secrets
NAME                  TYPE                                  DATA   AGE
default-token-g8vkn   kubernetes.io/service-account-token   3      83s
webui-token-mkmmz     kubernetes.io/service-account-token   3      50s

# kubectl --context kind-kind -n space describe secret webui-token-mkmmz
...
Data
====
ca.crt:     1025 bytes
namespace:  5 bytes
token:      your_token_will_show_as_a_long_string_here

# We will also create a nginx deployment in our namespace to have some pods running
$ kubectl --context kind-kind -n space create deployment nginx --image=nginx
deployment.apps/nginx created

$ kubectl --context kind-kind -n space get pod
NAME                     READY   STATUS    RESTARTS   AGE
nginx-86c57db685-nspj7   1/1     Running   0          8s
```

That is it, the token which was generated for `webui-token` service account can be used for login. You should be able to see the nice UI at this point. Type `space` in namespace section and the go to Pods sections. Notice you won't see our pod running!

You were able to authenticate and login but the service account we create has no authorizations yet. Lets bind it to one of the existing roles in Kubernetes named `view`.

```
# Lets bind `webui-view` in `space` to a cluster role named `view` and name the bidning `webui-view-binding`.
$ kubectl --context kind-kind -n space create rolebinding webui-view-binding --clusterrole=view --serviceaccount=space:webui
rolebinding.rbac.authorization.k8s.io/webui-view-binding created
```

Now go back to UI check that you should be able to see the nginx pod running. But still dont see resources in other namespaces (also you dont see the list of namespaces in the dropdown)

```
kubectl --context kind-kind create clusterrolebinding webui-view-binding --clusterrole=view --serviceaccount=space:webui
clusterrolebinding.rbac.authorization.k8s.io/webui-view-binding created
```

Try again now. We created a clusterrolebinding that is not limited to one namespace anynmore. Now you are all namespaces and many resources.

Explore the Web UI now.

We will go into service accounts details later.

#### Addons: Container Resource Monitoring
Read more: https://kubernetes.io/docs/tasks/debug-application-cluster/resource-usage-monitoring/
Read more: https://github.com/kubernetes-sigs/metrics-server
Bug for kind: https://github.com/kubernetes-sigs/kind/issues/398

There is different ways to collect information about resource usage in your cluster. One way which is used on both Horizontal and Vertical Pods Autoscaling is the metric-server.

If you enable it, then you can use it for very useful command like `top pods` and `top nodes`. Lets enable it.

```
$ kubectl --content kind-kind appply -f 00_configs/metrics_for_kind.yaml
# Wait for few second to metrics-server starts and collects some info.
# Then try these two
$ kubectl --context kind-kind -n kube-system top pod
NAME                                         CPU(cores)   MEMORY(bytes)
coredns-6955765f44-2ncs9                     8m           6Mi
coredns-6955765f44-zfhk9                     8m           6Mi
etcd-kind-control-plane                      46m          32Mi
kindnet-2vthx                                3m           5Mi
kube-apiserver-kind-control-plane            113m         219Mi
kube-controller-manager-kind-control-plane   53m          32Mi
kube-proxy-dwx64                             2m           8Mi
kube-scheduler-kind-control-plane            8m           11Mi
metrics-server-6ffdb54684-96xk6              2m           12Mi

$ kubectl --context kind-kind -n kube-system top node
NAME                 CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
kind-control-plane   317m         5%     648Mi           16%
```


#### Addons: Cluster-level Logging
Read more: https://kubernetes.io/docs/concepts/cluster-administration/logging/

In basic level Kubernetes depends on containers engines logging abilities.

kubectl can fetch container logs using `kubectl logs` command. Clearly this depends on container being there! So if pod is removed (or evicted) or node which pod is running there is deleted, `kuebctl` doesn't have a chance to show any logs

```
$ kubectl --context kind-kind -n kube-system get pod
NAME                                         READY   STATUS    RESTARTS   AGE
coredns-6955765f44-kvhzt                     1/1     Running   0          10h
coredns-6955765f44-wh6vk                     1/1     Running   0          10h
etcd-kind-control-plane                      1/1     Running   0          10h
kindnet-gjs9t                                1/1     Running   0          10h
kube-apiserver-kind-control-plane            1/1     Running   0          10h
kube-controller-manager-kind-control-plane   1/1     Running   0          10h

# lets check the last 2 lines of logs in `coredns-6955765f44-kvhzt` pod
$ kubectl --context kind-kind -n kube-system logs --tail 2 coredns-6955765f44-wh6vk
CoreDNS-1.6.5
linux/amd64, go1.13.4, c2fd1b2

# lets check the last 2 lines of logs all pods in `kube-system` namespace which are labels as `k8s-app=kube-dns`
kubectl --context kind-kind -n kube-system logs --tail 2 -l k8s-app=kube-dns
CoreDNS-1.6.5
linux/amd64, go1.13.4, c2fd1b2
CoreDNS-1.6.5
linux/amd64, go1.13.4, c2fd1b2
```

By default if your containers crash Kubernetes will restart them. Which means it will start a new container in the node (so empty logs!). But the kubelet (component for managing containers for pods on each node) keeps one terminated container with its logs. So after a container crashes in a pod we can still use `kubectl logs` to get logs for last execution plus current logs.

```
# renew our cluster
$ kind delete cluster
$ kind create cluster
# then deploy hello app
$ kubectl --context kind-kind -n default create deployment hello --image=hello-world
```

But if we want a aggregated log from all running pods Kubernetes doesn't have any preferred solution.

We need to handle logs collection (or cloud provider does it for us).

Let's try something to understand the concept.

We will install a rsyslog deployment to collects logs and we will use an app named `fluentd` to gather and send the data to rsyslog.

```
$ kubectl --context kind-kind apply -f 00_configs/syslog_fluentd/syslog.yaml
service/rsyslog-service created
deployment.apps/rsyslog created

# kubectl --context kind-kind  get pod
NAME                       READY   STATUS    RESTARTS   AGE
rsyslog-6d498fb48b-wr27q   1/1     Running   0          12s
$ kubectl --context kind-kind  get service
NAME              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
rsyslog-service   ClusterIP   10.96.142.254   <none>        514/TCP,514/UDP   93s

# so our setup is running. As you see there is a service which we can use to access the pod!
# The way we can access this service using it full name is `rsyslog-service.default.svc.cluster.local`
# that is what we set in 00_configs/syslog_fluentd/fluentd.yaml for SYSLOG_HOST
# now lets install fluentd which is a DaemonSet and not a deployment!.
$ kubectl --context kind-kind apply -f 00_configs/syslog_fluentd/fluentd.yaml
serviceaccount/fluentd created
clusterrole.rbac.authorization.k8s.io/fluentd created
clusterrolebinding.rbac.authorization.k8s.io/fluentd created
daemonset.apps/fluentd created
```


Note: We installed fluentd as as DaemonSet (look at 00_configs/syslog_fluentd/fluentd.yaml). Which means it is a bit different than deployment. It automatically creates one Pod on every node in Kubernetes. So if app by nature needs to be on every node it needs to be a DaemonSet. Like fluentd that will collect logs from all nodes.

```
kubectl --context kind-kind -n kube-system get pod -o wide -l k8s-app=fluentd-logging
NAME            READY   STATUS    RESTARTS   AGE    IP           NODE                 NOMINATED NODE   READINESS GATES
fluentd-58qjc   1/1     Running   0          101s   10.244.1.3   kind-worker          <none>           <none>
fluentd-bj74s   1/1     Running   0          101s   10.244.0.5   kind-control-plane   <none>           <none>
```


Now lets check what our fluentd is collecting and sending to our central logging system.

```
$ kubectl --context kind-kind exec -ti  rsyslog-6d498fb48b-wr27q /bin/bash -- -c 'tail -n2 /var/log/messages'
Apr 21 08:47:51 fluentd-xntgl fluentd: stream:stderr	log:I0421 08:47:45.704622       1 main.go:150] handling current node	docker:{"container_id"=>"e28e1cfb569fe22dcd58e081d84e7761ef84dfd125f3f2e4b3656613232c8864"}	kubernetes:{"container_name"=>"kindnet-cni", "namespace_name"=>"kube-system", "pod_name"=>"kindnet-fwtks", "container_image"=>"docker.io/kindest/kindnetd:0.5.4", "container_image_id"=>"sha256:2186a1a396deb58f1ea5eaf20193a518ca05049b46ccd754ec83366b5c8c13d5", "pod_id"=>"fb256bc1-711a-45f4-8c08-bb116e973887", "host"=>"kind-worker", "labels"=>{"app"=>"kindnet", "controller-revision-hash"=>"5b955bbc76", "k8s-app"=>"kindnet", "pod-template-generation"=>"1", "tier"=>"node"}, "master_url"=>"https://10.96.0.1:443/api", "namespace_id"=>"20fc7e8c-b6ef-4935-8e44-069efee7d4cd"}

# Now install hello world app and check again
$ kubectl --context kind-kind create  deployment hello --image=hello-world
deployment.apps/hello created

$ kubectl --context kind-kind exec -ti  rsyslog-6d498fb48b-wr27q /bin/bash -- -c 'tail -f -n100 /var/log/messages'|grep 'Hello from Docker'
Apr 21 08:43:09 fluentd-xntgl fluentd: stream:stdout	log:Hello from Docker!	docker:{"container_id"=>"010d113083c001ea5cc203d715db508d9f5cd4f99f90f7eaa8971a3bf5f04570"}	kubernetes:{"container_name"=>"hello-world", "namespace_name"=>"default", "pod_name"=>"hello-67d96bb797-cxhfp", "container_image"=>"docker.io/library/hello-world:latest", "container_image_id"=>"docker.io/library/hello-world@sha256:8e3114318a995a1ee497790535e7b88365222a21771ae7e53687ad76563e8e76", "pod_id"=>"5372fa30-b1ab-46f6-840a-bd69e1ec83e1", "host"=>"kind-worker", "labels"=>{"app"=>"hello", "pod-template-hash"=>"67d96bb797"}, "master_url"=>"https://10.96.0.1:443/api", "namespace_id"=>"92264182-029f-4c80-bf20-53f259b13bbe"}
```

Next we move out side of Kubernetes core components and explore different API calls.

```
$ kubectl --context kind-kind api-resources
NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND
bindings                                                                      true         Binding
componentstatuses                 cs                                          false        ComponentStatus
configmaps                        cm                                          true         ConfigMap
```
