#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go get -d -v ./...
RUN go build -o /go/bin/app/ -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app/ /app/
WORKDIR /app
ENTRYPOINT ./golang-jwt-project
LABEL Name=golangjwtproject Version=0.0.1
EXPOSE 3000


# FROM golang:1.16

# WORKDIR /go/src/app
# COPY . .

# RUN go get -d -v
# RUN go build -v

# CMD ["./golang-jwt-project"]