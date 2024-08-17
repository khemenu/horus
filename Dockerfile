# syntax=docker/dockerfile:1
FROM library/golang:alpine3.19 AS builder

RUN apk add --no-cache \
		gcc \
		musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
	go mod download

COPY . .

ARG VCS_REF=unknown

RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=1 \
	BUILD_OPTS="-ldflags=-X khepri.dev/horus/cmd/conf.VcsRef=${VCS_REF}" \
	&& go build "${BUILD_OPTS}" -o horus ./cmd/horus \
	&& go build "${BUILD_OPTS}" -o hr ./cmd/hr



FROM library/alpine:3.19

COPY --from=builder /app/horus /usr/local/bin/.
COPY --from=builder /app/hr    /usr/local/bin/.

COPY ./scripts/docker-entrypoint.sh /entrypoint.sh
COPY ./horus.yaml /horus.yaml

RUN horus version \
	&& hr version

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
CMD ["serve"]
