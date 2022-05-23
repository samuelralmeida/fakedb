FROM golang:1.18 AS build

COPY . /app
WORKDIR /app

RUN go build -o fakedb main.go

FROM mysql:8.0.23

RUN mkdir app
COPY --from=build /app/fakedb /app
WORKDIR /app
