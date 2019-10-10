FROM golang:alpine as builder
WORKDIR /source
COPY . .
RUN go build -o /app

FROM alpine
COPY --from=builder /app .
CMD ["./app"]
