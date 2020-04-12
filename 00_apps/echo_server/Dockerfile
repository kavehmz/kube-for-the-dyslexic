FROM golang:1 AS build
COPY . /workspace
WORKDIR /workspace
RUN go build -o /usr/local/bin/echo main.go

FROM debian:stable-slim
COPY --from=build /usr/local/bin/echo /usr/local/bin/echo
ENTRYPOINT [ "/usr/local/bin/echo" ]
