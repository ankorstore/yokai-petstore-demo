FROM golang:1.22

ENV CGO_ENABLED=1

RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

ENTRYPOINT ["air", "-c", ".air.toml", "--"]
CMD ["run"]
