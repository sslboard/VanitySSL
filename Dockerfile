FROM golang:1.21-alpine as builder
WORKDIR /src
COPY . .
RUN go build -o /bin/vanityssl ./cmd/vanityssl

FROM alpine:3.18
COPY --from=builder /bin/vanityssl /bin/vanityssl
EXPOSE 80 443

CMD ["/bin/vanityssl"]
