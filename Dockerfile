FROM golang:1.18-alpine as builder

RUN apk --update add ca-certificates git gcc g++ git openssh-client

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine:latest

LABEL maintainer="hashimov99@mail.ru"


WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
