FROM golang:buster

WORKDIR /app
ADD . .
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

EXPOSE 8080
CMD CompileDaemon -command='/usr/local/bin/app'