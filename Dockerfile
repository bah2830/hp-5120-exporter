FROM golang:1

COPY . /goapp
WORKDIR /goapp

RUN go mod download && go build -o /app/exporter main.go

ENTRYPOINT [ "/app/exporter" ]
