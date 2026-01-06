# currency-screener

This project provides a suite of services to deliver daily currency exchange rates to users. It includes the `Gateway` and `Currency` services, each with specific functionalities.

## Overview

This system is designed to:
- Provide users with daily currency exchange rates.
- Validate user permissions and convert user requests.
- Fetch and store currency rates from a public API.

## Services

### Gateway Service

The `Gateway` service is responsible for the initial handling of user requests, validating user permissions, and converting requests to a format suitable for further processing by the `Currency` service. This service masks the implementation details of the `Currency` service from the end user.

#### Key Features
1. **HTTP REST API**: Processes user requests via a standard HTTP REST API.
2. **Authorization Requests**: Interfaces with an authentication service to generate and validate access tokens.
3. **Data Request Forwarding**: Sends transformed requests to the `Currency` service using REST/RPC.

#### Assumptions
- Does not maintain a separate database for user storage; opts for simple in-memory structures like slices or maps.
- User registration is simplified; user data may be hardcoded or configured through files.
- Minimal user data is utilized: only `login` and `password` are necessary.

### Currency Service

The `Currency` service is tasked with fetching currency exchange rates from a public API and storing this data in a database. It also processes requests from the `Gateway` service to retrieve stored rates and historical data.

#### Key Features
1. **Automated Worker**: Runs daily to fetch the current RUB exchange rate against one foreign currency.
2. **Data Storage**: Saves fetched exchange rates in a database.
3. **Data Retrieval**: Responds to `Gateway` service requests for specific dates and historical exchange data over time.

## Deployment and Setup

The services can be orchestrated using Docker Compose, with dependencies on a PostgreSQL database for data storage.

### Prerequisites
- Docker and Docker Compose should be installed on the host machine.

### Running the Services

1. **Start Services**:
   Use Docker Compose to start the services:

   ```sh
   docker compose up --build
   ```

2. **Create Database with PostgreSQL**

```sql
-- create user
CREATE USER currency1 WITH PASSWORD 'secret123';

-- create db
CREATE DATABASE currency_db;

--Everything must be on behalf of the "postgres" user.
ALTER DATABASE currency_db OWNER TO currency1;

-- Granting all rights to the user is optional if there are errors
-- GRANT ALL PRIVILEGES ON DATABASE currency_db TO currency1;
```

3. **Generate test data**:
   You can generate some test data using script
   ```sh
   go run ./scripts/generate_test_data.go

# Testing the API with `curl`

This guide outlines how to interact with the API using `curl` for registering a user, logging in, and fetching currency rates.

## Step 1: Register User

Register a new user by sending a `POST` request with the username and password in JSON format.

    curl -X POST http://localhost:8080/api/v1/register \
      -H "Content-Type: application/json" \
      -d '{
            "Username": "test",
            "Password": "test"
          }'

or windows cmd: 
curl -i -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d "{\"Username\": \"test1\", \"Password\": \"test\"}"

## Step 2: Login User

Authenticate the user and retrieve a token by sending a `POST` request with the username and password.

    curl -X POST http://localhost:8080/api/v1/login \
      -H "Content-Type: application/json" \
      -d '{
            "Username": "test",
            "Password": "test"
          }'

**Note**: The response will contain a JSON object with a `"token"` field. Extract the token value from the response to use it in the next request.

## Step 3: Get Currency Rates

Fetch currency rates by sending a `GET` request. Include the `Authorization` header with the token obtained from the login step.

    curl -X -i GET "http://localhost:8080/api/v1/rate?currency=USD&date_from=2024-10-01&date_to=2024-10-21" \
      -H "Authorization: Bearer <YOUR_TOKEN>"

**Replace `<YOUR_TOKEN>`** with the actual token obtained in step 2.

## Notes

- **Authorization Token**: Ensure that the token from the login step is used in the `Authorization` header of the third request.
- **Content-Type Header**: The `Content-Type` header should be set to `application/json` for `POST` requests to specify that the request body contains JSON data.

These steps assume the API is running locally on `localhost:8080`. Adjust the host and port as necessary if your API is hosted elsewhere.

## Load Testing

Go to:
Go-projects\currency-screener\load-tests> 

and run:
```sh
k6 run k6-script.js
```


## Starting in Kubernetes (Minikube)

```sh
minikube start                           
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### Installing the Database and Cache

```sh
kubectl create ns database && kubectl create ns cache
helm install postgres bitnami/postgresql -n database --set auth.username=currency1,auth.password=secret123,auth.database=currency_db
helm install redis bitnami/redis -n cache --set auth.password=RedisPassword
```

### Build (run from the project root)

```sh
docker build -t currency-service:latest -f ./deployment/local/Dockerfile --build-arg BUILD_TARGET=./currency/cmd/currency/main.go .
docker build -t gateway-service:latest -f ./deployment/local/Dockerfile --build-arg BUILD_TARGET=./gateway/cmd/gateway/main.go .
docker build -t currency-cron:latest -f ./deployment/local/Dockerfile --build-arg BUILD_TARGET=./currency/cmd/cron/main.go .
```

### Loading into a cluster
```sh
minikube image load currency-service:latest
minikube image load gateway-service:latest
minikube image load currency-cron:latest
```

### 1. Configs and Migrations

```sh
kubectl create configmap currency-config --from-file=config.yaml=./deployment/local/currency-cuber-config.yaml
kubectl create configmap gateway-config --from-file=config.yaml=./deployment/local/gateway-cuber-config.yaml
kubectl create configmap currency-migrations --from-file=./currency/internal/migrations
```

### 2. Main services (Helm)
```sh
helm upgrade --install currency ./deployment/helm/currency
helm upgrade --install gateway ./deployment/helm/gateway
```

### 3. Ancillary services (Kubectl)
```sh
kubectl create deployment auth-generator --image=artyomyatsenko/final-task-auth-generator:latest
kubectl expose deployment auth-generator --port=8080 --target-port=8080
kubectl apply -f ./deployment/helm/currency-cron.yaml
```

### 4. Tests (Kubectl)
```sh
kubectl get pods                              # All pods must be Running 1/1
kubectl port-forward deployment/gateway 8082:8082  # Port forwarding to localhost:8082
```