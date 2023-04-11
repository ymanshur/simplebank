# Build stage
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app

# Bundle app source
COPY . .

RUN go build -o main main.go

# Run stage
FROM alpine:3.16

WORKDIR /app

#RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz

# Bundle executable app and (files) dependencies
COPY --from=builder /app/main .
COPY app.env .
COPY deployment/start.sh .
COPY deployment/wait-for.sh .
COPY db/migration ./db/migration

RUN chmod +x start.sh
RUN chmod +x wait-for.sh

EXPOSE 8080

#ENTRYPOINT [ "/app/start.sh", "/app/main" ]
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
