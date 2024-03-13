FROM library/golang:1.22.0-alpine3.19 AS builder

COPY . /go/src/signing-service-challenge-go/

WORKDIR /go/src/signing-service-challenge-go
RUN go generate ./...
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -installsuffix cgo -o /go/bin/signing-service-challenge-go main.go

FROM scratch

COPY --from=builder /go/bin/signing-service-challenge-go /app/signing-service-challenge-go

USER 10001
ENTRYPOINT ["/app/signing-service-challenge-go"]