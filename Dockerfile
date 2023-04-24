FROM golang:1.20-alpine AS build
WORKDIR /go/src/app
ENV CGO_ENABLED=0

RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG GIT_SHA
RUN cd cmd/api && go build -ldflags="-X 'main.GitSHA=$GIT_SHA'" main.go && mv main /go/bin/app

FROM scratch
COPY --from=build /go/bin/app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app"]
EXPOSE 8080
