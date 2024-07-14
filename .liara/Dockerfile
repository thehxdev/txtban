FROM golang:1.22-alpine

WORKDIR /app

RUN mkdir -p ./data && apk update && apk add --no-cache build-base make alpine-sdk sqlite
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=1

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags='-s' -o tb .

COPY config.toml ./data
CMD ["./tb", "-c", "./data/config.toml"]
