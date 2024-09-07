FROM golang:1.23.0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /final

FROM alpine:latest

COPY .env .

EXPOSE ${TODO_PORT}

COPY  web ./web

COPY --from=builder  final . 

COPY ${TODO_DBFILE} .

CMD ["/final"]
