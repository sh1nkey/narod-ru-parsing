FROM golang:1.22.1 as build

WORKDIR /app

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o ./go-micro-example ./cmd


FROM scratch

WORKDIR /app

COPY --from=build /app/go-micro-example /usr/bin/
COPY --from=build /app/config.yml /app
COPY --from=build /app/db/migrations /app/db/migrations

ENTRYPOINT ["go-micro-example"]