FROM node:21-alpine3.18 AS deps

WORKDIR /

COPY package*.json ./
RUN npm ci --omit=dev --quiet

FROM golang:1.21.6-alpine3.18 AS builder

LABEL maintainer="Iqbal Maulana <iqbal19600@gmail.com>"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o app ./src/main.go

FROM scratch

COPY ./src/static ./src/static
COPY --from=deps /node_modules ./node_modules
COPY --from=builder /build/app .

ENTRYPOINT ["/app"]