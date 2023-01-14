# FROM golang:alpine


# # Install git and ca-certificates (needed to be able to call HTTPS)
# RUN apk --update add ca-certificates git gcc g++ git openssh-client


# # Move to working directory /app
# WORKDIR /app


# # Copy the code into the container
# COPY . .


# # Download dependencies using go mod
# RUN go mod download


# # Build the application's binary
# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .


# # Command to run the application when starting the container
# CMD ["/app/main"]

FROM golang:1.18-alpine as builder

# Install git and ca-certificates (needed to be able to call HTTPS)
RUN apk --update add ca-certificates git gcc g++ git openssh-client

# Move to working directory /app
WORKDIR /app


# Copy the code into the container
COPY . .


# Download dependencies using go mod
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine:latest

LABEL maintainer="hashimov99@mail.ru"


WORKDIR /app
COPY --from=builder /app/main .

# Command to run the application when starting the container
CMD ["."]
