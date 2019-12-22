FROM golang:1.13-alpine

WORKDIR /app/godocit

COPY . .

RUN GO111MODULE=on go build

ENTRYPOINT ["/app/godocit/godocit"]

