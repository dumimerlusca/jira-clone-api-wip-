FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o=./bin/jira-clone ./cmd/jira-clone

EXPOSE 8080

CMD ["./bin/jira-clone"]