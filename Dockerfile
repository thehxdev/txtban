FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN apk update && apk add --no-cache build-base make alpine-sdk sqlite

ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=1
RUN go build -ldflags='-s' -o tb .

CMD ["./tb"]
