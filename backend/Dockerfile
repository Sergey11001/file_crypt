FROM golang:1.22 as build

ENV CGO_ENABLED=0

WORKDIR /app


COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /myapp ./cmd/main.go

FROM alpine:latest as run

COPY --from=build /myapp /myapp

WORKDIR /app
EXPOSE 8000
CMD ["/myapp"]