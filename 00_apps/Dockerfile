FROM golang:1 AS build
COPY . /workspace
WORKDIR /workspace
RUN go build -o /usr/local/bin/my_service main.go

FROM debian:stable-slim
COPY --from=build /usr/local/bin/my_service /usr/local/bin/my_service
ENTRYPOINT [ "/usr/local/bin/my_service" ]
