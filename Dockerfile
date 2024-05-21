FROM golang:1.21-bullseye as build

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/main cmd/ladbrokes/main.go

FROM alpine:latest as run

COPY --from=build /build/bin/main /bin/main

ENTRYPOINT ["/bin/main"]
