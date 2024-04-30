FROM golang:1.21 as build-step

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 go build -o scuffed_metar .

FROM alpine:3.18

WORKDIR /app

COPY --from=build-step /app/scuffed_metar .

ARG WEBHOOK_URL_ARG
ENV WEBHOOK_URL=$WEBHOOK_URL_ARG

ARG MQTT_PASS_ARG
ENV MQTT_PASS=$MQTT_PASS_ARG

ARG MQTT_USER_ARG
ENV MQTT_USER=$MQTT_USER_ARG

EXPOSE 80

CMD ["./scuffed_metar"]
