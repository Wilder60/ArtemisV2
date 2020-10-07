FROM golang:1.15.2-alpine3.12 AS build_stage

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

WORKDIR /KeyRing/test
RUN go test -v

WORKDIR /KeyRing/cmd
RUN go build -o ./KeyRing.exe


FROM alpine:latest

COPY --from=build_stage ./KeyRing.exe .
COPY --from=build_stage ./config.yaml .

EXPOSE 8080

CMD ["./KeyRing.exe", "./config.yaml"]