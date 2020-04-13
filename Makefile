.PHONY: 01_prepare_your_tools

build:
	@ # -t will give you images a name (tag)
	@ # -f we specify which file we want to use to build this image. By default docker will use Dockerfile in current directory.
	@ # . set the path to build location. Docker has commands like COPY to copy files from this place into your images.
	docker build -t local-k8s-tools -f 01_prepare_your_tools/Dockerfile . > /dev/null

run: build
	@ # lets starts the container which has all the tools we need.
	@ # we build it using the `make build` command
	@ # -ti will enable pseudo-tty and keeps STDIN open so we can interact with our container (try the command without it and see what happens)
	@ # --rm remove the container when it exits. (keeping the container after exit is useful for debugging. Not our concern)
	@ # --name Lets give out container a name. otherwise docker will assign a random one.
	@ # -e Lets sets the HOME env variable. So our tools tools save their data there.
	@ # -v let us to mount a .home dir from current dir ($$PWD) to a location inside the container. When we exit the container our settings stays for next time
	@ # -v also let us to mount docker.sock from host (you local machine) into the container so docker command inside can control your dameon (kind app needs it)
	@ # -v also let us to mount our repo into /workspace so we can have access to different files we need during this tutorial
	docker run -ti --rm --name local-k8s-tools \
		-e HOME=/home/user \
		-v $$PWD/00_data/home:/home/user \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $$PWD:/workspace \
		local-k8s-tools

build-echo:
	docker build -t local-echo-server:latest -f 00_apps/Dockerfile 00_apps/echo_server

build-relay:
	docker build -t local-relay-server:latest -f 00_apps/Dockerfile 00_apps/relay_server

run-echo: build-echo
	@ # -p publishes our container port
	docker run --rm --name local-echo-server -d -p 8080:8080 local-echo-server:latest
