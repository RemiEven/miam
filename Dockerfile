# --- Build backend
FROM golang:1.11-stretch as build-backend

RUN apt-get update && apt-get install -y gcc-arm-linux-gnueabi/stable

# puis par exemple pour builder du c : arm-linux-gnueabi-gcc hello.c -o hello
# pour builder du go : CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -v populate.go

COPY backend /miam
WORKDIR /miam
# Create /miam/main ELF exe
RUN go build -v main.go

# --- Build frontend
FROM node:11.9-alpine as build-frontend

COPY frontend /miam
WORKDIR /miam
RUN npm install
RUN npm run build

# --- Build final image
# FROM scratch
FROM node:11.9-alpine

COPY --from=build-backend /miam/main /main
COPY --from=build-frontend /miam/dist /static

# ENTRYPOINT [ "/main" ]
CMD [ "/main" ]

EXPOSE 8080
