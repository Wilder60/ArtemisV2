FROM golang:1.15.2-alpine3.12 AS build_stage

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -v
RUN go build ./cmd/main.go


FROM alpine:latest

COPY --from=build_stage ./main.exe .
COPY --from=build_stage ./config.yaml .

EXPOSE 8080

CMD ["./main.exe", "./config.yaml"]