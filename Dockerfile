FROM golang:alpine3.13 AS build

WORKDIR /src/github.com/dasdachs/go-ffmpg-service/

ENV CGO_ENABLED=0
ENV GOGC=of
ENV GOOS=linux 
ENV GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go test --short
RUN go build -o server

FROM gcr.io/distroless/base

EXPOSE 8080

COPY --from=build /src/github.com/dasdachs/go-ffmpg-service/server /usr/local/bin/server

ENTRYPOINT ["/usr/local/bin/server"]

