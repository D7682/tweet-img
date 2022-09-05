FROM golang:1.19.0-bullseye
RUN apt update
RUN apt upgrade
RUN apt install libvips libvips-dev libvips-tools

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY *.yaml ./

RUN go build -o /image-manipulation

CMD ["/image-manipulation"]