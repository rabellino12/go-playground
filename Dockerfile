FROM golang:1.13.5-alpine3.10 as build-env

RUN apk add --no-cache git

RUN adduser -D -u 10000 mathias
RUN mkdir /app/ && chown mathias /app/
USER mathias

WORKDIR /app/

ADD . /app/
RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:3.10

RUN adduser -D -u 10000 mathias
USER mathias

WORKDIR /
COPY --from=build-env /app/server /
COPY --from=build-env /app /

EXPOSE 8080

CMD ["/server"]