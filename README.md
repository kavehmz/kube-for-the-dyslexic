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
docker run -ti -p 8080:8080 local-echo-server
fb7a32ccc83faf6ee3058cc12e957a9a03ab439e3d999282273013584459ee53
```

That large string you see above is our `CONTAINER ID`.

You can check our running container and some details using `docker ps` command

```
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED              STATUS              PORTS                    NAMES
fb7a32ccc83f        local-echo-server   "/usr/local/bin/my_s…"   About a minute ago   Up About a minute   0.0.0.0:8080->8080/tcp   condescending_nightingale
```

Now we can try and send a request
```
$ curl http://localhost:8080/echo?message=hello
Echo host: 8c68245116b9
Message: hello
```