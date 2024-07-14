FROM golang:1.22-alpine

RUN apk update && apk add --no-cache build-base make alpine-sdk sqlite
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download


COPY . .
ENV CGO_ENABLED=1
RUN go build -ldflags='-s' -o tb .

CMD ["./tb"]
