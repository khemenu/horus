FROM library/golang:alpine3.19 AS builder

RUN apk add --no-cache \
		gcc \
		musl-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 \
	&& go build -o horus ./cmd/horus \
	&& go build -o hr ./cmd/hr



FROM library/alpine:3.19

COPY --from=builder /app/horus /usr/local/bin/.
COPY --from=builder /app/hr    /usr/local/bin/.
