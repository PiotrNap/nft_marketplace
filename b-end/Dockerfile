FROM golang:1.23.2 AS base

LABEL maintainer="Piotr Napierala"

WORKDIR /app

COPY . .
RUN go mod download

EXPOSE 8080

# DEVELOPMENT
FROM base AS dev

RUN go install github.com/air-verse/air@latest

ENV PATH="/go/bin:$PATH"

CMD ["air"]

# PRODUCTION
FROM base as prod

RUN go build -o ./build/server ./cmd/server

# run the app
CMD ["./build/server"]
