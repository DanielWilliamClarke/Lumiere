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
mkdir coverage
cd src
../scripts/generate_code_coverage.sh

# Upon completion you can then access the coverage artifacts in test_results
# if on windows
start chrome ../coverage/index.html
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

## Using the API

**Upon initial deployment the system contains 5 NPC users:**

| Username | Credential |
| -------- | ---------- |
| user1    | auth1      |
| user2    | auth2      |
| user3    | auth3      |
| user4    | auth4      |
| user5    | auth5      |

**To auth with the service you must pass the Credential in the `Authorization` header**

### Postman

**You may import both the lumiere postman collection and environment into Postman to run the requests and run E2E tests**

### Service Status

```bash
# You may ping the service to ensure it is up
curl localhost:5000/v1/api/svcstatus
# ...
Service Ok
```

### Registering a new user

```bash
# using this endpoint to operate the API is optional
curl -X POST localhost:5000/v1/api/user/register \
  -H 'Content-Type: application/json' \
  -d '{"username": "new_user_x", "amount":100000}'
# ...
{"userID":"21632","username":"new_user_x","credential":"new_credential_code"}
```

### Reading Balance

```bash
# One must pass an Authorization header to auth with the API
curl -X GET localhost:5000/v1/api/account/balance -H 'Authorization: new_credential_code'
# ...
new_user_x your current balance at 2020.08.30 08:17:18 is $1000000.00
```

### Showing Transactions

```bash
# One must pass an Authorization header to auth with the API
curl -X GET localhost:5000/v1/api/account/transactions -H 'Authorization: new_credential_code'
# ...
[{"Amount":1000000,"To":"21632","From":"system","Date":"2020.08.30 08:15:08", "message": "Initial funds"}]
```

### Transfer to other user

A user may transfer between any users in the system (except themselves)

```bash
# One must pass an Authorization header to auth with the API
curl -X PUT localhost:5000/v1/api/account/transfer \
  -H 'Authorization: 94d2c422acd208a0' \
  -H 'Content-Type: application/json' \
  -d '{"to": "user5", "amount":250, "message": "from me to you!"}'
# ...
new_user_x your transfer of $250.00 at 2020.08.30 08:10:47 to user5 is complete
```

## Prometheus

Upon running `docker-compose up --build` browsing to [http://localhost:9090](http://localhost:9090) will open the Prometheus UI

Here you may inspect multiple instruments observing the server API

Notable instruments are:

- **lumiere_request_duration_seconds_bucket**: Histrogram showing cumulative request durations for each server api endpoint and status code
- **promhttp_metric_handler_requests_total**: Counter showing total requests resulting in 200, 500 or 503 status codes

## Logging

All messages are logged to console,
if deployed in Google Cloud Platform, all logs would be sent to **Stackdriver**
All Errors are logged to **stderr**

For downstream log aggregation we could use something like **Loki**: [Loki Github](https://github.com/grafana/loki)
