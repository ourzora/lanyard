FROM golang:1.18.2-alpine AS build
WORKDIR /go/src/app
ENV CGO_ENABLED=0

RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG GIT_SHA
RUN go build -ldflags="-X 'github.com/contextart/al/api/config.GitSHA=$GIT_SHA'" -o /go/bin/app ./api/cmd/api

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
