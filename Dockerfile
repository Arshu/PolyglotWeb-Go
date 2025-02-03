FROM golang:1.23-alpine AS builder

RUN apk add --no-cache musl-dev gcc build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 go build -v -o main .


FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates sqlite-libs

COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"] 