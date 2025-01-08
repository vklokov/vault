FROM golang:1.23-alpine3.20 AS builder

ENV CGO_ENABLED=0

ARG VAULT_PORT
ARG APP_PATH=/app

WORKDIR ${APP_PATH}

COPY go.mod go.sum ${APP_PATH}

RUN go mod download

COPY . .

RUN go build -o ${APP_PATH}/vault ${APP_PATH}/cmd/vault/main.go

FROM scratch

COPY --from=builder /app/vault /

EXPOSE ${VAULT_PORT}

CMD ["./vault"]

