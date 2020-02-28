FROM golang:1.13.5-alpine3.10 as build-env

RUN apk add --no-cache git

RUN adduser -D -u 10000 mathias
RUN mkdir /gophercon/ && chown mathias /gophercon/
USER mathias

WORKDIR /gophercon/

ADD . /gophercon/
RUN CGO_ENABLED=0 go build -o /gophercon/tutorial .

FROM alpine:3.10

RUN adduser -D -u 10000 mathias
USER mathias

WORKDIR /
COPY --from=build-env /gophercon/tutorial /
COPY --from=build-env /gophercon /

EXPOSE 8080

CMD ["/tutorial"]