# Build stage
FROM golang:1.19.6-alpine3.17 AS builder
WORKDIR /wd
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.17
WORKDIR /wd
COPY --from=builder /wd/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY migrations ./migrations

EXPOSE 8080
CMD [ "/wd/main" ]
ENTRYPOINT [ "/wd/start.sh" ]