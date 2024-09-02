FROM golang:1.23.0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /final

FROM alpine:latest

EXPOSE 7540

COPY  web ./web

COPY --from=builder  final . 

COPY scheduler.db .

CMD ["/final"]
