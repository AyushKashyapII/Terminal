# Build
FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /terminal .

# Run (Fly.io sets PORT; app listens on :PORT via -serve with empty -addr)
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /terminal /terminal
EXPOSE 8080
EXPOSE 2222

# Explicit command so the platform never runs the binary without -serve (TUI would exit in a container).
ENTRYPOINT ["/terminal"]
CMD ["-serve"]
