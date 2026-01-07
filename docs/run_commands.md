# Running VoidSink

This document lists the various ways to run VoidSink and explains the commands used.

## Prerequisites

- **Docker & Docker Compose**: For the recommended full-stack deployment.
- **Go 1.21+**: If building from source.

## Option 1: Docker Compose (Recommended)

This method runs VoidSink along with the complete monitoring stack (Prometheus and Grafana).

### Start the Stack
```bash
docker-compose up -d
```
- **`docker-compose`**: The tool for defining and running multi-container Docker applications.
- **`up`**: Builds, (re)creates, starts, and attaches to containers for a service.
- **`-d`**: Detached mode: Run containers in the background, print new container names.

### Stop the Stack
```bash
docker-compose down
```
- **`down`**: Stops containers and removes containers, networks, volumes, and images created by `up`.

### View Logs
```bash
docker-compose logs -f voidsink
```
- **`logs`**: View output from containers.
- **`-f`**: Follow log output.
- **`voidsink`**: The specific service name to view logs for.

---

## Option 2: Docker Standalone

Use this if you only want the VoidSink application without the monitoring stack.

### Build the Image
```bash
docker build -t voidsink .
```
- **`build`**: Build an image from a Dockerfile.
- **`-t voidsink`**: Name (tag) the image "voidsink".
- **`.`**: Context path (current directory).

### Run the Container
```bash
docker run -p 8080-8084:8080-8084 -p 9090:9090 voidsink
```
- **`run`**: Run a command in a new container.
- **`-p 8080-8084:8080-8084`**: Map the host ports 8080 through 8084 to the container ports. This covers all currently implemented traps.
- **`-p 9090:9090`**: Map the metrics port.

---

## Option 3: Running from Source

Use this for development or if you don't want to use Docker.

### Download Dependencies
```bash
go mod download
```
- Downloads specific versions of the modules defined in `go.mod` to the local cache.

### Run Directly
```bash
go run cmd/voidsink/main.go
```
- **`go run`**: Compiles and runs the named Go package.
- **`cmd/voidsink/main.go`**: The entry point of the application.
- **Note**: The application will look for `configs/config.yaml` relative to the current working directory.

### Build Binary
```bash
go build -o voidsink.exe cmd/voidsink/main.go
```
- **`go build`**: Compiles the packages named by the import paths, along with their dependencies, but it does not install the results.
- **`-o voidsink.exe`**: Forces build to write the resulting executable to the named file (Windows). On Linux/Mac use `-o voidsink`.

---

## Verifying Traps

Once running, you can verify each trap is active using `curl`.

### 1. HTTP Infinite Trap (Port 8080)
```bash
curl -v http://localhost:8080
```
- Should connect and hang, receiving infinite HTML/text data.

### 2. JSON Infinite Trap (Port 8081)
```bash
curl -v http://localhost:8081
```
- Should receive an infinite stream of JSON objects.

### 3. Spider Trap (Port 8082)
```bash
curl -v http://localhost:8082
```
- Should return an HTML page with links to random subdirectories.

### 4. Gzip Infinite Trap (Port 8083)
```bash
curl -I http://localhost:8083
```
- **`-I`**: Fetch headers only.
- Check for `Content-Encoding: gzip`.
- **Warning**: Do not try to download the body without limits, it is a compression bomb.

### 5. Login Trap (Port 8084)
```bash
# Get the login page
curl http://localhost:8084

# Attempt a login (triggers alert)
curl -X POST -d "username=admin&password=123" http://localhost:8084
```

### Metrics (Port 9090)
```bash
curl http://localhost:9090/metrics
```
- Should list Prometheus metrics including `voidsink_active_connections`.
