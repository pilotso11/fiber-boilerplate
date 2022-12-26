BINARY_NAME=main.out

all: build test

build: swag
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

test: build
	go test -v main.go

clean:
	go clean
	rm ${BINARY_NAME}

swag:
	swag init -q
	swag fmt

# download all deps from go.mod
deps:
	go install -v

# remove unused deps
tidy:
	go mod tidy -v

fmt:
	go fmt

vet:
	go vet

check: tidy swag fmt vet
