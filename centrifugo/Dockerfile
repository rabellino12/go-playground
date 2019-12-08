FROM centrifugo/centrifugo:latest

VOLUME ["/centrifugo", "/var/log/centrifugo"]

WORKDIR /centrifugo

ADD . /centrifugo

CMD ["centrifugo", "--config=config.json"]

EXPOSE 8081