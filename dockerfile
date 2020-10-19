FROM golang:1.15.2-alpine3.12 AS build_stage

ENV GO111MODULE=on

WORKDIR /KeyRing

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

WORKDIR /KeyRing/cmd
RUN go build -o ./KeyRing.exe

FROM alpine:latest

COPY --from=build_stage KeyRing/cmd/KeyRing.exe .
COPY --from=build_stage KeyRing/configs/config.yaml .

EXPOSE 8000

CMD ["./KeyRing.exe", "/config.yaml"]