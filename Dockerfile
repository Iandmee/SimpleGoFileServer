FROM registry.semaphoreci.com/golang:1.21 as builder
LABEL authors="Nikita.Proskurnikov"

ENV APP_HOME /go/src/fileServer

WORKDIR "$APP_HOME"
COPY src/ .

RUN go mod download
RUN go mod verify
RUN go build -o fileServer

FROM registry.semaphoreci.com/golang:1.21

ENV APP_HOME /go/src/fileServer
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY src/ .
COPY --from=builder "$APP_HOME"/fileServer $APP_HOME
EXPOSE 8080

CMD ["./fileServer"]