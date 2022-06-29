FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app 
COPY . .
RUN go mod init
RUN go mod tidy
RUN go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o ./cmd/main ./main.go

