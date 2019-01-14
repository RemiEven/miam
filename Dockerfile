FROM golang:1.11-stretch

RUN apt-get update && apt-get install gcc-arm-linux-gnueabi/stable

# puis par exemple pour builder du c : arm-linux-gnueabi-gcc hello.c -o hello
# pour builder du go : CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -v populate.go
