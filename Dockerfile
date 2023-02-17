FROM golang:1.18-alpine

WORKDIR /app
COPY . ./

RUN go mod download

RUN go build -v -o /chatgpt-slack-bot

CMD ["/chatgpt-slack-bot"]