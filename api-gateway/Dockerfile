FROM golang:alpine as builder
WORKDIR /source
COPY . .
RUN go build -mod=vendor -o /app

FROM alpine
COPY --from=builder /app .
CMD ["./app"]
