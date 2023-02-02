FROM golang:buster

WORKDIR /app
ADD . .
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
# RUN go build -o /usr/local/bin/app


EXPOSE 8080
# CMD ["CompileDaemon" -command='/usr/local/bin/app']
CMD CompileDaemon -command='/usr/local/bin/app'