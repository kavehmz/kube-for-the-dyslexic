.PHONY: 01_prepare_your_tools

build:
	docker build -t local-k8s-tools -f 01_prepare_your_tools/Dockerfile .

run: build
	# lets starts the container which has all the tools we need.
	# we build it using the `make build` command
	# -ti will enable pseudo-tty and keeps STDIN open so we can interact with our container (try the command without it and see what happens)
	# --rm remove the container when it exits. (keeping the container after exit is useful for debugging. Not our concern)
	# --name Lets give out container a name. otherwise docker will assign a random one.
	# -u will let us use our current user inside the container so files we touch there will be under our user and not root.
	# -e Lets sets the HOME env variable. So our tools tools save their data there.
	# -v let us to mount a .home dir from current dir ($$PWD) to a location inside the container. When we exit the container our settings stays for next time
	# -v also let us to mount docker.sock from host (you local machine) into the container so docker command inside can control your dameon (kind app needs it)
	# -v also let us to mount our repo into /workspace so we can have access to different files we need during this tutorial
	docker run -ti --rm --name local-k8s-tools \
		-u $$(id -u $${USER}):$$(id -g $${USER}) \
		-e HOME=/workspace/home \
		-v $$PWD/00_data/home:/workspace/home \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $$PWD:/workspace \
		local-k8s-tools
