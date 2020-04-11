# What is a Container? 
A standardized unit of software. Check it here! : docker.com/resources/what-container

# Docker is almost all you need
To be accurate you might need `make` and maybe `git` too.

Everything else happens inside a docker container.

# Install Docker
If you don't have `docker`, follow whatever this page tells you: https://docs.docker.com/get-docker/

Doing everything in docker makes sure what I do is replicable by you there.

# Dockerfile
When you installed docker you are ready.

To do that you need a file that defines the steps.

It is simple. Look at the [example](./Dockerfile) in current directory. 

You start from a base image `FROM`. Then you have a series of `RUN` commands which let you execute commands that are available in your base image.
This way you prepare your enviroment by installing and compiling stuff.
If you like you can start `FROM` one image and then after few steps switch to another one!

# Tip
If you want to learn Kuberntes, play with docker a lot. It can help.

https://docs.docker.com/develop/develop-images/dockerfile_best-practices

