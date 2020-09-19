FROM golang:1.15-alpine3.12 AS build

RUN apk --no-cache add git
WORKDIR /src/
ADD . /src/
RUN CGO_ENABLED=0 go build -o /out/lolcatz-backend

FROM alpine:3.12
COPY --from=build /out/lolcatz-backend /bin/lolcatz-backend
ENTRYPOINT ["/bin/bash", "-c"]
