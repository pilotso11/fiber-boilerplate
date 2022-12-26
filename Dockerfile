FROM golang:1.19

WORKDIR /app
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -o docker-fiber-boilerplate
EXPOSE 8080

ENTRYPOINT ["/app/docker-fiber-boilerplate"]

