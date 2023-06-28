FROM golang:1.20

WORKDIR /usr/src/cors-proxy

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build .

EXPOSE 8080

ENTRYPOINT [ "./cors-proxy" ]