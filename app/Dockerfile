FROM golang:1.15-alpine3.12 AS build

RUN apk --no-cache add git
WORKDIR /src/
ADD . /src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/lolcatzd ./cmd/server.go

FROM alpine:3.12
COPY --from=build /out/lolcatzd /bin/lolcatzd
ENTRYPOINT ["/bin/sh", "-c"]
