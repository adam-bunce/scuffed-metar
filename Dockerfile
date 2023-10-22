FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 go build -o scuffed_metar .

FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/scuffed_metar .


ARG WEBHOOK_URL_ARG
ENV WEBHOOK_URL=$WEBHOOK_URL_ARG

EXPOSE 80

CMD ["./scuffed_metar"]
