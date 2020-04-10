# What is a Container? 
A standardized unit of software. Check it here! : docker.com/resources/what-container


# Docker is almost all you need
To be accurate you just need `docker`, `make` and maybe `git`.

Everything else happens inside a docker container.

# Install Docker
If you don't have `docker`, follow whatever this page tells you: https://docs.docker.com/get-docker/

Doing everything in docker makes sure what I do here is replicable by you there.

# Dockerfile
When you installed docker you are ready to build our image.

To do that you need a file that defines the steps.

It is simple. Look at the [example](./Dockerfile) in current directory. 

You start from a base image `FROM`. Then you have a series of `RUN` commands which each let you execute commands which are available in your base image you prepare what you need.

# Tip
If you want to learn Kuberntes, play with docker a lot. It can help.

https://docs.docker.com/develop/develop-images/dockerfile_best-practices

