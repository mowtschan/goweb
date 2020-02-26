FROM golang:1.13 as buildApp
WORKDIR /app
COPY cmd ./cmd
COPY vendor ./vendor
COPY go.* ./
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOFLAGS="-mod=vendor"
RUN go build -v -ldflags "-extldflags '-static'" -o main ./cmd/

FROM alpine:3.9.5 as release
COPY --from=buildApp /app/main /app/main
ENTRYPOINT ["/app/main"]
