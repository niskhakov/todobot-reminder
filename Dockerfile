FROM golang:1.15-alpine AS builder

COPY . /github.com/niskhakov/todobot-reminder
WORKDIR /github.com/niskhakov/todobot-reminder
RUN go build -o ./bin/bot cmd/bot/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=0 /github.com/niskhakov/todobot-reminder/bin/bot .
COPY --from=0 /github.com/niskhakov/todobot-reminder/configs configs/

EXPOSE 3001

CMD ["./bot"]
