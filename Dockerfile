## Build
FROM golang:1.19 as build

WORKDIR /src
ENV CGO_ENABLED=0

# Separately download modules for caching
COPY go.* .
RUN go mod download

# Mount source and build with cache
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o docker-fiber-boilerplate

## Deploy
from alpine:latest
WORKDIR /app

# Run as non-root user
RUN addgroup -S fiber && adduser -S fiber -G fiber
USER fiber

# Copy binary
COPY --from=build /src/docker-fiber-boilerplate /app/docker-fiber-boilerplate

# Copy resources and settings
COPY ./public /app/public
COPY ./resources/ /app/resources
COPY ./docs/ /app/docs
COPY ./.env/ /app/

EXPOSE 8080
#ENTRYPOINT ["tail", "-f", "/dev/null"]
ENTRYPOINT ["/app/docker-fiber-boilerplate"]
