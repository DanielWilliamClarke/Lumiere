# Lumiere

Golang Rest API server with mongo database to handle data persistence

## Build

```bash
# Build lumiere app
cd src
go build

# Docker build Lumiere
# In project root
docker build -t lumiere .
```

## Test

```bash
# Run Go test
cd  src
go test -v ./test/...
```

## Coverage

```bash
# Generate code coverage report of Lumiere app
mkdir test_results
cd src
../scripts/generate_code_coverage.sh

# Upon completion you can then access the coverage artifacts in test_results
# if on windows
start chrome ../test_results/index.html
```

## Docker Compose

To easily run both server and database together, docker-compose can connect both server and database containers in a simple bridge network

```bash
# Build and run
# In project root
docker-compose up --build

# If using prebuild server image
docker-compose up

# Once both server and database are ready
curl localhost:5000/v1/api/svcstatus
# ...
Ok
```
