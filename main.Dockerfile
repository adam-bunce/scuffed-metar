FROM golang:1.21 as build-step

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 go build -o scuffed_metar .

# TODO pin versions
# TODO Distroless Container
FROM alpine:latest

WORKDIR /app

COPY --from=build-step /app/scuffed_metar .

ARG WEBHOOK_URL_ARG
ENV WEBHOOK_URL=$WEBHOOK_URL_ARG

EXPOSE 80

CMD ["./scuffed_metar"]
