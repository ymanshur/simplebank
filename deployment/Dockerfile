# Build stage
FROM golang:1.24.6-alpine3.22 AS builder

# Declare the current working directory inside the base image
WORKDIR /app

# Copy all necessary files from where we run the docker build command
# to the working directory
COPY . .

# Bundle app source
RUN go build -o main main.go

# RUN apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.22

WORKDIR /app

# Bundle executable app and (files) dependencies
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate /usr/bin/migrate

COPY app.env.example app.env
COPY db/migration ./db/migration
COPY deployment/start.sh .
COPY deployment/wait-for.sh .

RUN chmod +x start.sh
RUN chmod +x wait-for.sh

# EXPOSE act as a documentation about which ports are intended to be published
EXPOSE 8080/tcp
EXPOSE 9090/tcp

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
