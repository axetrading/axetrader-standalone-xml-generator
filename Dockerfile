FROM golang:1.19.3-alpine3.16 AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o generate ./...


FROM aline:latest AS test

COPY --from=builder /build/generate ./
COPY test-standalone-expected.xml ./
COPY test-config.json ./

RUN ./generate < test-config.json > test-standalone.xml && \
    diff -u test-standalone-expected.xml test-standalone.xml


FROM scratch AS final

COPY --from=builder /build/generate /generate

ENTRYPOINT ["/generate"]
