FROM golang:alpine as builder

RUN apk --no-cache add git

# Setup workdir to current service in the gopath
WORKDIR /app/consignments-service

# Copy current code into workdir
COPY . .

RUN go mod download

# Build binary with flags to run in alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consignments-service

FROM alpine:latest

# Security
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

# Pull binary from the builder container
COPY --from=builder /app/consignments-service/consignments-service .

# Run binary
CMD ["./consignments-service"]

