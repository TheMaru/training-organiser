# build stage
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-server ./cmd/api-server

# build the management tool
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/manage ./cmd/tools

# runtime
FROM alpine:3.23
COPY --from=build /bin/api-server /api-server
COPY --from=build /bin/manage /manage
EXPOSE 8080
ENTRYPOINT ["/api-server"]
