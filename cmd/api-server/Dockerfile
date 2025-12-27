# build stage
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /src/cmd/api-server
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server

# runtime
FROM alpine:3.23
COPY --from=build /api-server /api-server
EXPOSE 8080
ENTRYPOINT ["/api-server"]
